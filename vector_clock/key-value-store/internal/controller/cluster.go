package controller

import (
	"fmt"
	"hash"
	"hash/fnv"
	"log"
	"vectory_clock/key-value-store/internal/hashring"
	"vectory_clock/pkg/model"
)

// INode is an interface all cluster nodes implement for use in the consistent hash ring.
type INode interface {
	hashring.ICacheNode
	GetValue(k string) (*model.ValueWithClock, error)
	SetValueWithClock(key string, v *model.ValueWithClock) (*model.ValueWithClock, error)
}

// Cluster aggregates nodes and routing logic for reads/writes.
type Cluster struct {
	hashRingObj *hashring.HashRing
	config      *ClusterConfig
}

// ClusterConfig holds cluster-wide, operator-tunable parameters.
type ClusterConfig struct {
	readQuorum    int
	writeQuorum   int
	totalReplicas int
	virtualNodes  int
	hashFunction  func() hash.Hash64
}

type ClusterOption func(*ClusterConfig) *ClusterConfig

func WithHashFunction(f func() hash.Hash64) ClusterOption {
	return func(cfg *ClusterConfig) *ClusterConfig { cfg.hashFunction = f; return cfg }
}
func WithVirtualNodes(count int) ClusterOption {
	return func(cfg *ClusterConfig) *ClusterConfig { cfg.virtualNodes = count; return cfg }
}
func WithReadQuorum(q int) ClusterOption {
	return func(cfg *ClusterConfig) *ClusterConfig { cfg.readQuorum = q; return cfg }
}
func WithWriteQuorum(q int) ClusterOption {
	return func(cfg *ClusterConfig) *ClusterConfig { cfg.writeQuorum = q; return cfg }
}
func WithTotalReplicas(r int) ClusterOption {
	return func(cfg *ClusterConfig) *ClusterConfig { cfg.totalReplicas = r; return cfg }
}

// NewCluster builds a Cluster configured from options; hashes use FNV-1a by default.
func NewCluster(opts ...ClusterOption) (*Cluster, error) {
	defaultConfig := &ClusterConfig{
		readQuorum: 2, writeQuorum: 2, totalReplicas: 3, virtualNodes: 3,
		hashFunction: fnv.New64a,
	}
	for _, opt := range opts {
		defaultConfig = opt(defaultConfig)
	}
	// TIP: Strict quorums for safety (R + W > N)
	if defaultConfig.readQuorum <= 0 || defaultConfig.writeQuorum <= 0 ||
		defaultConfig.totalReplicas <= 0 ||
		defaultConfig.readQuorum > defaultConfig.totalReplicas ||
		defaultConfig.writeQuorum > defaultConfig.totalReplicas ||
		defaultConfig.readQuorum+defaultConfig.writeQuorum < defaultConfig.totalReplicas {
		return nil, fmt.Errorf("invalid config: R=%d, W=%d, N=%d",
			defaultConfig.readQuorum, defaultConfig.writeQuorum, defaultConfig.totalReplicas)
	}
	return &Cluster{
		config: defaultConfig,
		hashRingObj: hashring.InitHashRing(
			hashring.SetVirtualNodes(defaultConfig.virtualNodes),
			hashring.SetReplicationFactor(defaultConfig.totalReplicas),
			hashring.EnableVerboseLogs(true),
			hashring.SetHashFunction(defaultConfig.hashFunction),
		),
	}, nil
}

// AddNode registers a new node in the cluster.
func (c *Cluster) AddNode(node INode) error {
	if err := c.hashRingObj.AddNode(node); err != nil {
		log.Printf("[ERROR] failed to add node %s: %v", node.GetIdentifier(), err)
		return fmt.Errorf("failed to add node %s: %w", node.GetIdentifier(), err)
	}
	log.Printf("[INFO] Node %s added to hash ring", node.GetIdentifier())
	return nil
}

// RemoveNode removes a node from the cluster.
func (c *Cluster) RemoveNode(node INode) error {
	if err := c.hashRingObj.RemoveNode(node); err != nil {
		log.Printf("[ERROR] failed to remove node %s: %v", node.GetIdentifier(), err)
		return fmt.Errorf("failed to remove node %s: %w", node.GetIdentifier(), err)
	}
	log.Printf("[INFO] Node %s removed from hash ring", node.GetIdentifier())
	return nil
}

// Get performs a quorum get on the key (R nodes); resolves conflicts if needed.
func (c *Cluster) Get(k string) (*model.ValueWithClock, error) {
	nodes, err := c.hashRingObj.GetNodesForKey(k)
	if err != nil {
		log.Printf("[ERROR] hash ring get failed: %v", err)
		return nil, fmt.Errorf("failed to get values for key %s: %w", k, err)
	}
	if len(nodes) == 0 {
		return nil, fmt.Errorf("no nodes available for key %s", k)
	}
	values := make([]*model.ValueWithClock, 0, c.config.readQuorum)
	nodesSlice := make([]INode, 0, c.config.readQuorum)
	for _, node := range nodes {
		value, err := node.(INode).GetValue(k)
		if err != nil {
			log.Printf("[WARN] Could not get key=%s from node=%s: %v", k, node.GetIdentifier(), err)
			continue
		}
		nodesSlice = append(nodesSlice, node.(INode))
		values = append(values, value)
		if len(values) >= c.config.readQuorum {
			break
		}
	}
	return c.resolveConflicts(nodesSlice, k, values), nil
}

// resolveConflicts reconciles potentially divergent values using LWW.
// If repair is needed, sets consistent value across nodes.
func (c *Cluster) resolveConflicts(nodes []INode, k string, values []*model.ValueWithClock) *model.ValueWithClock {
	if len(values) == 0 {
		return nil
	}
	latest := values[0]
	for i := 1; i < len(values); i++ {
		comp := latest.Clock.Compare(values[i].Clock)
		if comp == -1 {
			// values[0] is older; values[i] wins (LWW)
			log.Printf("[REPAIR] Key=%s Detected older value, updating primary", k)
			latest.Value = values[i].Value
			latest.Clock = latest.Clock.Merge(values[i].Clock)
			c.setValueOnNode(nodes[0], k, latest)
		} else if comp == 1 {
			// values[i] is older; repair that replica
			log.Printf("[REPAIR] Key=%s Repair back-propagate newer value to stale replica", k)
			values[i].Value = latest.Value
			values[i].Clock = latest.Clock.Merge(values[i].Clock)
			c.setValueOnNode(nodes[i], k, values[i])
		}
	}
	return latest
}

// setValueOnNode forces a specific key/value on a node.
func (c *Cluster) setValueOnNode(node INode, k string, v *model.ValueWithClock) (*model.ValueWithClock, error) {
	log.Printf("[DEBUG] SyncRepair SET: key=%s node=%s value=%v clock=%v", k, node.GetIdentifier(), v.Value, v.Clock)
	return node.SetValueWithClock(k, v)
}

// Set writes the value to W nodes (quorum); async replicates to the rest.
func (c *Cluster) Set(k string, v *model.ValueWithClock) (*model.ValueWithClock, error) {
	nodes, err := c.hashRingObj.GetNodesForKey(k)
	if err != nil || len(nodes) == 0 {
		return nil, fmt.Errorf("failed to get nodes for key %s: %w", k, err)
	}
	count := 0
	var lastValue *model.ValueWithClock
	for _, node := range nodes {
		n := node.(INode)
		if count >= c.config.writeQuorum {
			go c.setValueOnNode(n, k, v)
		} else {
			val, err := c.setValueOnNode(n, k, v)
			if err != nil {
				log.Printf("[ERROR] Failed to set value on node=%s: %v", n.GetIdentifier(), err)
				continue
			}
			count++
			lastValue = val
		}
	}
	return lastValue, nil
}
