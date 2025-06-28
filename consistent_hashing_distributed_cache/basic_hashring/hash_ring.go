package basichashring

import (
	"errors"
	"fmt"
	"hash"
	"hash/fnv"
	"log"
	"sort"
	"sync"
)

var (
	ErrNoConnectedNodes = errors.New("no connected nodes available")
	ErrNodeExists       = errors.New("node already exists")
	ErrNodeNotFound     = errors.New("node not found")
	ErrInHashingKey     = errors.New("error in hashing the key")
)

type ICacheNode interface {
	GetIdentifier() string
}

type hashRingConfig struct {
	HashFunction func() hash.Hash64
	EnableLogs   bool
}

type HashRingConfigFn func(*hashRingConfig)

func SetHashFunction(f func() hash.Hash64) HashRingConfigFn {
	return func(config *hashRingConfig) {
		config.HashFunction = f
	}
}

func EnableVerboseLogs(enabled bool) HashRingConfigFn {
	return func(config *hashRingConfig) {
		config.EnableLogs = enabled
	}
}

type HashRing struct {
	mu         sync.RWMutex
	config     hashRingConfig
	nodes      sync.Map
	sortedKeys []uint64
}

func InitHashRing(opts ...HashRingConfigFn) *HashRing {
	config := &hashRingConfig{
		HashFunction: fnv.New64a,
		EnableLogs:   false,
	}
	for _, opt := range opts {
		opt(config)
	}
	return &HashRing{
		config:     *config,
		sortedKeys: make([]uint64, 0),
	}
}

func (ring *HashRing) AddNode(node ICacheNode) error {
	ring.mu.Lock()
	defer ring.mu.Unlock()

	hashValue, err := ring.generateHash(node.GetIdentifier())
	if err != nil {
		return fmt.Errorf("%w: %s", ErrInHashingKey, node.GetIdentifier())
	}

	if _, exists := ring.nodes.Load(hashValue); exists {
		return fmt.Errorf("%w: node %s", ErrNodeExists, node.GetIdentifier())
	}

	ring.nodes.Store(hashValue, node)
	ring.sortedKeys = append(ring.sortedKeys, hashValue)

	// Maintain sorted order of hash keys for binary search
	sort.Slice(ring.sortedKeys, func(i, j int) bool { return ring.sortedKeys[i] < ring.sortedKeys[j] })

	if ring.config.EnableLogs {
		log.Printf("[HashRing] ➕ Added node: %s (hash: %d)", node.GetIdentifier(), hashValue)
	}
	return nil
}

func (ring *HashRing) RemoveNode(node ICacheNode) error {
	ring.mu.Lock()
	defer ring.mu.Unlock()

	hashValue, err := ring.generateHash(node.GetIdentifier())
	if err != nil {
		return fmt.Errorf("%w: %s", ErrInHashingKey, node.GetIdentifier())
	}

	if _, found := ring.nodes.LoadAndDelete(hashValue); !found {
		return fmt.Errorf("%w: %s", ErrNodeNotFound, node.GetIdentifier())
	}

	index, err := ring.search(hashValue)
	if err != nil {
		return err
	}

	ring.sortedKeys = append(ring.sortedKeys[:index], ring.sortedKeys[index+1:]...)

	if ring.config.EnableLogs {
		log.Printf("[HashRing] ❌ Removed node: %s (hash: %d)", node.GetIdentifier(), hashValue)
	}
	return nil
}

func (ring *HashRing) Get(key string) (ICacheNode, error) {
	ring.mu.RLock()
	defer ring.mu.RUnlock()

	hashValue, err := ring.generateHash(key)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInHashingKey, key)
	}

	index, err := ring.search(hashValue)
	if err != nil {
		return nil, err
	}

	nodeHash := ring.sortedKeys[index]
	if node, ok := ring.nodes.Load(nodeHash); ok {
		if ring.config.EnableLogs {
			log.Printf("[HashRing] 🔍 Key '%s' (hash: %d) mapped to node (hash: %d)", key, hashValue, nodeHash)
		}
		return node.(ICacheNode), nil
	}

	return nil, fmt.Errorf("%w: no node found for key %s", ErrNodeNotFound, key)
}

func (ring *HashRing) search(key uint64) (int, error) {
	if len(ring.sortedKeys) == 0 {
		return -1, ErrNoConnectedNodes
	}

	index := sort.Search(len(ring.sortedKeys), func(i int) bool {
		return ring.sortedKeys[i] >= key
	})

	// Wrap around the ring
	if index == len(ring.sortedKeys) {
		index = 0
	}
	return index, nil
}

// generateHash converts a string key to a uint64 hash using the configured hash function
func (ring *HashRing) generateHash(key string) (uint64, error) {
	h := ring.config.HashFunction()
	if _, err := h.Write([]byte(key)); err != nil {
		return 0, err
	}
	return h.Sum64(), nil
}
