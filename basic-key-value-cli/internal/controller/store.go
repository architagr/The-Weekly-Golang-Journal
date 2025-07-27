package controller

import (
	"errors"
	"sync"
	"time"
)

type item struct {
	value     string
	expiresAt time.Time
}

type Store struct {
	data map[string]item
	mu   sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		data: make(map[string]item),
	}
}

func (s *Store) Put(key, value string, ttlSeconds int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	expiration := time.Now().Add(time.Duration(ttlSeconds) * time.Second)
	s.data[key] = item{value: value, expiresAt: expiration}
}

func (s *Store) Get(key string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	it, exists := s.data[key]
	if !exists || time.Now().After(it.expiresAt) {
		// Key expired, or not found, try to delete it
		s.mu.RUnlock()
		s.mu.Lock()
		delete(s.data, key)
		s.mu.Unlock()
		s.mu.RLock()
		return "", errors.New("key not found or expired")
	}

	return it.value, nil
}

func (s *Store) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.data[key]; !exists {
		return errors.New("key not found")
	}
	delete(s.data, key)
	return nil
}
