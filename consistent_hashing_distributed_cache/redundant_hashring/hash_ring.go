package redundanthashring

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
	ErrNoNodesAvailable = errors.New("no connected nodes available")
	ErrNodeExists       = errors.New("node already exists")
	ErrNodeNotFound     = errors.New("node not found")
	ErrHashingKey       = errors.New("failed to hash key")
)

type ICacheNode interface {
	GetIdentifier() string
}

type hashRingConfig struct {
	VirtualNodes      int
	ReplicationFactor int
	HashFunction      func() hash.Hash64
	EnableLogs        bool
}

type HashRingConfigFn func(*hashRingConfig)

func SetVirtualNodes(count int) HashRingConfigFn {
	return func(cfg *hashRingConfig) {
		cfg.VirtualNodes = count
	}
}

func SetReplicationFactor(replication int) HashRingConfigFn {
	return func(cfg *hashRingConfig) {
		cfg.ReplicationFactor = replication
	}
}

func SetHashFunction(f func() hash.Hash64) HashRingConfigFn {
	return func(cfg *hashRingConfig) {
		cfg.HashFunction = f
	}
}

func EnableVerboseLogs(b bool) HashRingConfigFn {
	return func(cfg *hashRingConfig) {
		cfg.EnableLogs = b
	}
}

type HashRing struct {
	mu         sync.RWMutex
	config     hashRingConfig
	vNodeMap   sync.Map // hash → node
	hostSet    sync.Map // nodeID → bool
	sortedKeys []uint64
}

func InitHashRing(opts ...HashRingConfigFn) *HashRing {
	cfg := &hashRingConfig{
		VirtualNodes:      3,
		ReplicationFactor: 2,
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
			log.Printf("🧩 Virtual node added %s → %d", vID, h)
		}
	}
	ring.hostSet.Store(id, true)
	slices.Sort(ring.sortedKeys)
	return nil
}

func (ring *HashRing) RemoveNode(node ICacheNode) error {
	ring.mu.Lock()
	defer ring.mu.Unlock()

	id := node.GetIdentifier()
	if _, ok := ring.hostSet.Load(id); !ok {
		return ErrNodeNotFound
	}
	ring.hostSet.Delete(id)

	// remove all virtual nodes
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
	return nil
}

// ✅ GetPrimaryNode returns just one node (like V1 & V2)
func (ring *HashRing) GetPrimaryNode(key string) (ICacheNode, error) {
	ring.mu.RLock()
	defer ring.mu.RUnlock()

	h, err := ring.generateHash(key)
	if err != nil {
		return nil, err
	}

	idx := ring.search(h)
	node, _ := ring.vNodeMap.Load(ring.sortedKeys[idx])
	return node.(ICacheNode), nil
}

// ✅ GetNodesForKey returns N unique physical nodes for redundancy
func (ring *HashRing) GetNodesForKey(key string) ([]ICacheNode, error) {
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
	nodes := make([]ICacheNode, 0, ring.config.ReplicationFactor)

	start := ring.search(h)
	i := start

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
			nodes = append(nodes, n)
		}
		i++
		if i-start > len(ring.sortedKeys) {
			break // avoid infinite loop if not enough unique nodes
		}
	}

	if len(nodes) == 0 {
		return nil, ErrNoNodesAvailable
	}
	return nodes, nil
}

// 🧠 search returns index of first key ≥ hash or wraps around
func (ring *HashRing) search(h uint64) int {
	idx := sort.Search(len(ring.sortedKeys), func(i int) bool {
		return ring.sortedKeys[i] >= h
	})
	if idx == len(ring.sortedKeys) {
		return 0
	}
	return idx
}

func (ring *HashRing) generateHash(key string) (uint64, error) {
	h := ring.config.HashFunction()
	if _, err := h.Write([]byte(key)); err != nil {
		return 0, err
	}
	return h.Sum64(), nil
}
