package replicationhashring

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

var (
	ErrNoConnectedNodes = errors.New("no connected nodes available")
	ErrNodeExists       = errors.New("node already exists")
	ErrNodeNotFound     = errors.New("node not found")
	ErrHashingKey       = errors.New("failed to hash the key")
)

type ICacheNode interface {
	GetIdentifier() string
}

type hashRingConfig struct {
	VirtualNodes int
	HashFunction func() hash.Hash64
	EnableLogs   bool
}

type HashRingConfigFn func(*hashRingConfig)

func SetVirtualNodes(count int) HashRingConfigFn {
	return func(cfg *hashRingConfig) {
		cfg.VirtualNodes = count
	}
}

func SetHashFunction(fn func() hash.Hash64) HashRingConfigFn {
	return func(cfg *hashRingConfig) {
		cfg.HashFunction = fn
	}
}

func EnableVerboseLogs(enabled bool) HashRingConfigFn {
	return func(cfg *hashRingConfig) {
		cfg.EnableLogs = enabled
	}
}

type HashRing struct {
	mu         sync.RWMutex
	config     hashRingConfig
	hostMap    sync.Map // nodeID → timeAdded (for presence check)
	vNodeMap   sync.Map // hash → node
	sortedKeys []uint64 // sorted hash values (including virtual nodes)
}

func InitHashRing(opts ...HashRingConfigFn) *HashRing {
	cfg := &hashRingConfig{
		HashFunction: fnv.New64a,
		VirtualNodes: 3,
	}
	for _, opt := range opts {
		opt(cfg)
	}
	return &HashRing{
		config:     *cfg,
		sortedKeys: make([]uint64, 0),
	}
}

func (ring *HashRing) AddNode(node ICacheNode) error {
	ring.mu.Lock()
	defer ring.mu.Unlock()

	nodeID := node.GetIdentifier()
	if _, exists := ring.hostMap.Load(nodeID); exists {
		return fmt.Errorf("%w: %s", ErrNodeExists, nodeID)
	}

	virtualKeys := make([]uint64, 0, ring.config.VirtualNodes)
	for i := 0; i < ring.config.VirtualNodes; i++ {
		vNodeID := fmt.Sprintf("%s_%d", nodeID, i)
		h, err := ring.generateHash(vNodeID)
		if err != nil {
			return fmt.Errorf("%w for virtual node %s", ErrHashingKey, vNodeID)
		}
		ring.vNodeMap.Store(h, node)
		virtualKeys = append(virtualKeys, h)

		if ring.config.EnableLogs {
			log.Printf("[HashRing] ➕ Added virtual node %s → hash %d", vNodeID, h)
		}
	}

	ring.hostMap.Store(nodeID, struct{}{})
	ring.sortedKeys = append(ring.sortedKeys, virtualKeys...)
	slices.Sort(ring.sortedKeys)

	if ring.config.EnableLogs {
		log.Printf("[HashRing] ✅ Node %s added with %d virtual nodes", nodeID, ring.config.VirtualNodes)
	}
	return nil
}

func (ring *HashRing) RemoveNode(node ICacheNode) error {
	ring.mu.Lock()
	defer ring.mu.Unlock()

	nodeID := node.GetIdentifier()
	if _, exists := ring.hostMap.Load(nodeID); !exists {
		return fmt.Errorf("%w: %s", ErrNodeNotFound, nodeID)
	}

	for i := 0; i < ring.config.VirtualNodes; i++ {
		vNodeID := fmt.Sprintf("%s_%d", nodeID, i)
		h, err := ring.generateHash(vNodeID)
		if err != nil {
			return fmt.Errorf("%w for virtual node %s", ErrHashingKey, vNodeID)
		}
		ring.vNodeMap.Delete(h)

		// Remove from sortedKeys
		index := slices.Index(ring.sortedKeys, h)
		if index >= 0 {
			ring.sortedKeys = append(ring.sortedKeys[:index], ring.sortedKeys[index+1:]...)
		}

		if ring.config.EnableLogs {
			log.Printf("[HashRing] ❌ Removed virtual node %s (hash: %d)", vNodeID, h)
		}
	}
	ring.hostMap.Delete(nodeID)
	return nil
}

func (ring *HashRing) Get(key string) (ICacheNode, error) {
	ring.mu.RLock()
	defer ring.mu.RUnlock()

	h, err := ring.generateHash(key)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrHashingKey, key)
	}

	index, err := ring.search(h)
	if err != nil {
		return nil, err
	}

	nodeHash := ring.sortedKeys[index]
	if node, ok := ring.vNodeMap.Load(nodeHash); ok {
		if ring.config.EnableLogs {
			log.Printf("[HashRing] 🔍 Key '%s' (hash: %d) → node (hash: %d)", key, h, nodeHash)
		}
		return node.(ICacheNode), nil
	}

	return nil, fmt.Errorf("%w: key %s", ErrNodeNotFound, key)
}

func (ring *HashRing) search(hashValue uint64) (int, error) {
	if len(ring.sortedKeys) == 0 {
		return -1, ErrNoConnectedNodes
	}
	index := sort.Search(len(ring.sortedKeys), func(i int) bool {
		return ring.sortedKeys[i] >= hashValue
	})
	if index == len(ring.sortedKeys) {
		index = 0 // wrap around
	}
	return index, nil
}

func (ring *HashRing) generateHash(key string) (uint64, error) {
	h := ring.config.HashFunction()
	_, err := h.Write([]byte(key))
	if err != nil {
		return 0, err
	}
	return h.Sum64(), nil
}
