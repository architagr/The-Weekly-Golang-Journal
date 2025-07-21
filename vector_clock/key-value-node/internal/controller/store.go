package controller

import (
	"log"
	"sync"
	"vectory_clock/key-value-node/internal/config"
	"vectory_clock/pkg/model"
)

type Store struct {
	data sync.Map
	mu   sync.RWMutex
}

func NewStore() *Store {
	// Initialize the store with the provided parameters
	// This could include setting up connections to the key-value store, etc.
	return &Store{
		data: sync.Map{},
		mu:   sync.RWMutex{},
	}
}

func (s *Store) Get(key string) (*model.ValueWithClock, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	log.Println("Getting value from store:", key)
	value, ok := s.data.Load(key)
	if !ok {
		return nil, false
	}
	return value.(*model.ValueWithClock), true
}
func (s *Store) Set(key string, value *model.ValueWithClock) *model.ValueWithClock {
	s.mu.Lock()
	defer s.mu.Unlock()

	log.Println("Setting value in store:", key, value)
	if len(value.Clock) == 0 {
		s.mu.Unlock()
		v, ok := s.Get(key) // Load existing value to get the clock
		s.mu.Lock()
		if ok {
			log.Println("Found existing value for key:", key, v)
			v.Clock.Increment(config.NodeId) // Increment the clock for the current node
		} else {
			log.Println("No existing value found for key:", key)
			v = &model.ValueWithClock{
				Value: value.Value,
				Clock: model.VectorClock{config.NodeId: 0},
			}
		}
		log.Println("Setting new value with clock:", v.Clock)
		value.Clock = v.Clock.Copy()
	}
	log.Println("Storing value in map:", key, value)
	s.data.Store(key, value)
	return value
}
