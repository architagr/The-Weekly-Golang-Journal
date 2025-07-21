package controller

import (
	"fmt"
	"hash"
	"hash/fnv"
	"log"
	"vectory_clock/key-value-store/internal/hashring"
	"vectory_clock/pkg/model"
)

type INode interface {
	hashring.ICacheNode
	GetValue(k string) (*model.ValueWithClock, error)
	SetValueWithClock(key string, v *model.ValueWithClock) (*model.ValueWithClock, error)
}
type Cluster struct {
	hashRingObj *hashring.HashRing
	config      *ClusterConfig
}

type ClusterConfig struct {
	readQuorm     int
	writeQuorm    int
	totalReplicas int
	virtualNodes  int
	hashFunction  func() hash.Hash64
}

type ClusterOption func(*ClusterConfig) *ClusterConfig

func WithHashFunction(f func() hash.Hash64) ClusterOption {
	return func(cfg *ClusterConfig) *ClusterConfig {
		cfg.hashFunction = f
		return cfg
	}
}

func WithVirtualNodes(count int) ClusterOption {
	return func(cfg *ClusterConfig) *ClusterConfig {
		cfg.virtualNodes = count
		return cfg
	}
}

func WithReadQuorum(q int) ClusterOption {
	return func(cfg *ClusterConfig) *ClusterConfig {
		cfg.readQuorm = q
		return cfg
	}
}

func WithWriteQuorum(q int) ClusterOption {
	return func(cfg *ClusterConfig) *ClusterConfig {
		cfg.writeQuorm = q
		return cfg
	}
}

func WithTotalReplicas(r int) ClusterOption {
	return func(cfg *ClusterConfig) *ClusterConfig {
		cfg.totalReplicas = r
		return cfg
	}
}

func NewCluster(opts ...ClusterOption) (*Cluster, error) {
	defaultConfig := &ClusterConfig{
		readQuorm:     2,
		writeQuorm:    2,
		totalReplicas: 3,
		virtualNodes:  3,
		hashFunction:  fnv.New64a,
	}
	for _, opt := range opts {
		defaultConfig = opt(defaultConfig)
	}

	if defaultConfig.readQuorm <= 0 || defaultConfig.writeQuorm <= 0 ||
		defaultConfig.totalReplicas <= 0 || defaultConfig.readQuorm > defaultConfig.totalReplicas ||
		defaultConfig.writeQuorm > defaultConfig.totalReplicas || defaultConfig.readQuorm+defaultConfig.writeQuorm < defaultConfig.totalReplicas {
		return nil, fmt.Errorf("invalid cluster configuration: readQuorum=%d, writeQuorum=%d, totalReplicas=%d",
			defaultConfig.readQuorm, defaultConfig.writeQuorm, defaultConfig.totalReplicas)
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

func (c *Cluster) AddNode(node INode) error {
	if err := c.hashRingObj.AddNode(node); err != nil {
		return fmt.Errorf("failed to add node %s: %w", node.GetIdentifier(), err)
	}
	return nil
}

func (c *Cluster) RemoveNode(node INode) error {
	if err := c.hashRingObj.RemoveNode(node); err != nil {
		return fmt.Errorf("failed to remove node %s: %w", node.GetIdentifier(), err)
	}
	return nil
}

func (c *Cluster) Get(k string) (*model.ValueWithClock, error) {
	// Implementation for getting a value from the cluster
	// This will involve using the hash ring to find the appropriate node
	// and then fetching the value from that node.

	nodes, err := c.hashRingObj.GetNodesForKey(k)
	if err != nil {
		return nil, fmt.Errorf("failed to get values for key %s: %w", k, err)
	}
	if len(nodes) == 0 {
		return nil, fmt.Errorf("no nodes available for key %s", k)
	}
	values := make([]*model.ValueWithClock, 0, c.config.readQuorm)
	nodesSlice := make([]INode, 0, c.config.readQuorm)
	for _, node := range nodes {
		value, err := node.(INode).GetValue(k)
		if err != nil {
			return nil, fmt.Errorf("failed to get value from node %s: %w", node.GetIdentifier(), err)
		}
		nodesSlice = append(nodesSlice, node.(INode))
		values = append(values, value)
		if len(values) >= c.config.readQuorm {
			break
		}
	}
	return c.resolveConflicts(nodesSlice, k, values), nil
}

func (c *Cluster) resolveConflicts(nodes []INode, k string, values []*model.ValueWithClock) *model.ValueWithClock {
	// Implementation for resolving conflicts among the values fetched from different nodes
	// This could involve using vector clocks or other conflict resolution strategies.
	if len(values) == 0 {
		return nil
	}
	for i := 1; i < len(values); i++ {
		comp := values[0].Clock.Compare(values[i].Clock)
		if comp == 0 {
			continue // No conflict
		}
		if comp == -1 {
			// values[0] happened-before values[i], so we can keep values[i], as we are following Last write wins
			// This means values[0] is older, so we can discard it

			values[0].Value = values[i].Value
			values[0].Clock = values[0].Clock.Merge(values[i].Clock) // Copy the clock from the newer value
			c.setValueOnNode(nodes[0], k, values[0])                 // Update the node with the newer value
			continue
		}
		if comp == 1 {
			// values[0] happened-after values[i], so we can keep values[0]
			// This means values[i] is older, so we can discard it
			values[i].Value = values[0].Value
			values[i].Clock = values[i].Clock.Merge(values[0].Clock) // Discard the older value
			c.setValueOnNode(nodes[i], k, values[i])
			continue
		}
	}

	return values[0]
}

func (c *Cluster) setValueOnNode(node INode, k string, v *model.ValueWithClock) (*model.ValueWithClock, error) {
	// Implementation for setting a value on a specific node
	// This will involve calling the SetValueWithClock method on the node.
	log.Printf("Setting value for key %s on node %s", k, node.GetIdentifier())
	return node.SetValueWithClock(k, v)
}

func (c *Cluster) Set(k string, v *model.ValueWithClock) (*model.ValueWithClock, error) {
	// Implementation for setting a value in the cluster
	// This will involve using the hash ring to find the appropriate nodes
	// and then setting the value in those nodes.
	nodes, err := c.hashRingObj.GetNodesForKey(k)
	if err != nil {
		return nil, fmt.Errorf("failed to get nodes for key %s: %w", k, err)
	}
	if len(nodes) == 0 {
		return nil, fmt.Errorf("no nodes available for key %s", k)
	}
	count := 0
	var lastValue *model.ValueWithClock
	for _, node := range nodes {

		if count >= c.config.writeQuorm {
			go c.setValueOnNode(node.(INode), k, v)
		} else {
			v, err = c.setValueOnNode(node.(INode), k, v)
			if err != nil {
				return nil, fmt.Errorf("failed to set value on node %s: %w", node.GetIdentifier(), err)
			}
			count++
			lastValue = v // Last write wins
		}

	}

	return lastValue, nil
}
