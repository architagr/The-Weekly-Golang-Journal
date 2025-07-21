package controller

import (
	"log"
	"sync"
	"vectory_clock/key-value-node/internal/config"
	"vectory_clock/pkg/model"
)

// Store manages key-value pairs and their vector clocks for a node.
type Store struct {
	data sync.Map     // Thread-safe storage for all key-value entries
	mu   sync.RWMutex // Additional lock for complex read-write operations
}

// NewStore initializes a fresh key-value store.
func NewStore() *Store {
	return &Store{}
}

// Get returns the value with vector clock for the specified key, if present.
func (s *Store) Get(key string) (*model.ValueWithClock, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	log.Printf("[DEBUG] GET key=%s", key)
	value, ok := s.data.Load(key)
	if !ok {
		log.Printf("[DEBUG] GET key=%s not found", key)
		return nil, false
	}
	return value.(*model.ValueWithClock), true
}

// Set stores or updates a key-value entry along with its vector clock.
// Clocks are incremented for this node if not already set.
func (s *Store) Set(key string, value *model.ValueWithClock) *model.ValueWithClock {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(value.Clock) == 0 {
		// Try to load existing clock state
		s.mu.Unlock()
		v, ok := s.Get(key)
		s.mu.Lock()
		if ok {
			log.Printf("[DEBUG] SET key=%s incrementing existing vector clock %v", key, v.Clock)
			v.Clock.Increment(config.NodeId)
		} else {
			log.Printf("[DEBUG] SET key=%s creating new vector clock for node=%s", key, config.NodeId)
			v = &model.ValueWithClock{
				Value: value.Value,
				Clock: model.VectorClock{config.NodeId: 0},
			}
		}
		value.Clock = v.Clock.Copy()
	}
	log.Printf("[DEBUG] SET key=%s value=%v vectorClock=%v", key, value.Value, value.Clock)
	s.data.Store(key, value)
	return value
}
