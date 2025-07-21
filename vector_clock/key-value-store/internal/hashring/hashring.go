package hashring

import (
	"errors"
	"fmt"
	"hash"
	"hash/fnv"
	"log"
	"slices"
	"sort"
	"sync"
)

const (
	defaultVirtualNodes      = 3
	defaultReplicationFactor = 2
)

var (
	ErrNoNodesAvailable = errors.New("no connected nodes available")
	ErrNodeExists       = errors.New("node already exists")
	ErrNodeNotFound     = errors.New("node not found")
	ErrHashingKey       = errors.New("failed to hash key")
)

// ICacheNode is the host abstraction (for both physical and virtual/replicated nodes).
type ICacheNode interface {
	GetIdentifier() string
}

// Configuration for the hash ring; used for dependency injection and tuning.
type hashRingConfig struct {
	VirtualNodes      int
	ReplicationFactor int
	HashFunction      func() hash.Hash64
	EnableLogs        bool
}

// HashRingConfigFn is a functional option for customizing the hash ring.
type HashRingConfigFn func(*hashRingConfig)

func SetVirtualNodes(count int) HashRingConfigFn {
	return func(cfg *hashRingConfig) { cfg.VirtualNodes = count }
}
func SetReplicationFactor(replication int) HashRingConfigFn {
	return func(cfg *hashRingConfig) { cfg.ReplicationFactor = replication }
}
func SetHashFunction(f func() hash.Hash64) HashRingConfigFn {
	return func(cfg *hashRingConfig) { cfg.HashFunction = f }
}
func EnableVerboseLogs(b bool) HashRingConfigFn {
	return func(cfg *hashRingConfig) { cfg.EnableLogs = b }
}

// HashRing is a consistent hashing ring as used in backend clusters.
type HashRing struct {
	mu         sync.RWMutex
	config     hashRingConfig
	vNodeMap   sync.Map // hash → node
	hostSet    sync.Map // nodeID → bool
	sortedKeys []uint64 // sorted hash ring
}

// InitHashRing sets up a new hash ring with the given configuration.
func InitHashRing(opts ...HashRingConfigFn) *HashRing {
	cfg := &hashRingConfig{
		VirtualNodes:      defaultVirtualNodes,
		ReplicationFactor: defaultReplicationFactor,
		HashFunction:      fnv.New64a,
	}
	for _, opt := range opts {
		opt(cfg)
	}
	return &HashRing{
		config:     *cfg,
		sortedKeys: make([]uint64, 0),
	}
}

// AddNode builds virtual nodes for each hash slot and adds them to the ring.
func (ring *HashRing) AddNode(node ICacheNode) error {
	ring.mu.Lock()
	defer ring.mu.Unlock()
	id := node.GetIdentifier()
	if _, exists := ring.hostSet.Load(id); exists {
		return ErrNodeExists
	}

	for i := 0; i < ring.config.VirtualNodes; i++ {
		vID := fmt.Sprintf("%s#%d", id, i)
		h, err := ring.generateHash(vID)
		if err != nil {
			return fmt.Errorf("%w: %s", ErrHashingKey, vID)
		}
		ring.vNodeMap.Store(h, node)
		ring.sortedKeys = append(ring.sortedKeys, h)

		if ring.config.EnableLogs {
			log.Printf("[RING] Added virtual node %s → %d", vID, h)
		}
	}
	ring.hostSet.Store(id, true)
	slices.Sort(ring.sortedKeys)
	if ring.config.EnableLogs {
		log.Printf("[RING] Node %s fully added; ring now has %d unique hashes", id, len(ring.sortedKeys))
	}
	return nil
}

// RemoveNode deletes all vnodes for a host from the ring.
func (ring *HashRing) RemoveNode(node ICacheNode) error {
	ring.mu.Lock()
	defer ring.mu.Unlock()
	id := node.GetIdentifier()
	if _, ok := ring.hostSet.Load(id); !ok {
		return ErrNodeNotFound
	}
	ring.hostSet.Delete(id)

	// Remove all virtual nodes
	newKeys := make([]uint64, 0, len(ring.sortedKeys))
	for _, h := range ring.sortedKeys {
		val, ok := ring.vNodeMap.Load(h)
		if ok && val.(ICacheNode).GetIdentifier() == id {
			ring.vNodeMap.Delete(h)
			continue
		}
		newKeys = append(newKeys, h)
	}
	ring.sortedKeys = newKeys
	if ring.config.EnableLogs {
		log.Printf("[RING] Node %s removed. Ring now: %d keys", id, len(ring.sortedKeys))
	}
	return nil
}

// GetPrimaryNode returns the node responsible for the given key (like old V1 & V2).
func (ring *HashRing) GetPrimaryNode(key string) (ICacheNode, error) {
	ring.mu.RLock()
	defer ring.mu.RUnlock()

	h, err := ring.generateHash(key)
	if err != nil {
		return nil, err
	}
	if len(ring.sortedKeys) == 0 {
		return nil, ErrNoNodesAvailable
	}
	idx := ring.search(h)
	node, _ := ring.vNodeMap.Load(ring.sortedKeys[idx])
	return node.(ICacheNode), nil
}

// GetNodesForKey returns up to N unique physical nodes for redundancy (replicas).
func (ring *HashRing) GetNodesForKey(key string) (map[string]ICacheNode, error) {
	ring.mu.RLock()
	defer ring.mu.RUnlock()

	if len(ring.sortedKeys) == 0 {
		return nil, ErrNoNodesAvailable
	}

	h, err := ring.generateHash(key)
	if err != nil {
		return nil, err
	}

	seen := make(map[string]struct{})
	nodes := make(map[string]ICacheNode, ring.config.ReplicationFactor)

	start := ring.search(h)
	i := start
	// TIP: A ring walk across virtual nodes to gather N physical nodes (avoid duplicates).
	for len(nodes) < ring.config.ReplicationFactor {
		vHash := ring.sortedKeys[i%len(ring.sortedKeys)]
		node, ok := ring.vNodeMap.Load(vHash)
		if !ok {
			i++
			continue
		}
		n := node.(ICacheNode)
		id := n.GetIdentifier()
		if _, already := seen[id]; !already {
			seen[id] = struct{}{}
			nodes[n.GetIdentifier()] = n
		}
		i++
		if i-start > len(ring.sortedKeys) {
			break // infinite loop safety: not enough unique hosts in ring
		}
	}
	if len(nodes) == 0 {
		return nil, ErrNoNodesAvailable
	}
	return nodes, nil
}

// search is a ring binary-search: returns index where hash ≥ h or wraps around.
func (ring *HashRing) search(h uint64) int {
	idx := sort.Search(len(ring.sortedKeys), func(i int) bool {
		return ring.sortedKeys[i] >= h
	})
	if idx == len(ring.sortedKeys) {
		return 0
	}
	return idx
}

// generateHash hashes a string to a uint64 for ring use.
func (ring *HashRing) generateHash(key string) (uint64, error) {
	h := ring.config.HashFunction()
	if _, err := h.Write([]byte(key)); err != nil {
		return 0, err
	}
	return h.Sum64(), nil
}
