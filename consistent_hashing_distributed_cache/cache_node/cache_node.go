package cachenode

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

var (
	ErrNodeNotConnected = errors.New("node is not connected")
	ErrKeyNotFound      = errors.New("key not found")
)

// CacheNode represents a single in-memory cache server.
// In this simulated setup, each node can go down randomly after some time.
type CacheNode struct {
	isConnected bool     // Simulates node connection status
	Identifier  string   // Unique identifier for the node
	data        sync.Map // In-memory key-value storage
}

// InitCacheNode initializes a new cache node that will randomly go offline once.
// This is used to simulate real-world node failure in a distributed system.
func InitCacheNode(identifier string) *CacheNode {
	node := &CacheNode{
		Identifier:  identifier,
		isConnected: true,
	}
	node.simulateRandomFailure()
	return node
}

// simulateRandomFailure simulates a real-world flaky node by disconnecting
// the node once after a random delay between 0 to 10 seconds.
//
// This is intentional — the node *won’t* reconnect on its own.
// It’s meant to test how the hash ring handles unreachable nodes.
func (node *CacheNode) simulateRandomFailure() {
	delay := time.Duration(rand.Intn(10)) * time.Second

	log.Printf("[Node: %s] Simulating failure in %v...\n", node.Identifier, delay)
	time.AfterFunc(delay, func() {
		node.isConnected = false
		log.Printf("[Node: %s] Node disconnected.\n", node.Identifier)
	})
}

// GetIdentifier returns the node's identifier.
func (node *CacheNode) GetIdentifier() string {
	return node.Identifier
}

// ping checks if the node is connected. Returns an error if not.
func (node *CacheNode) ping() error {
	if !node.isConnected {
		return ErrNodeNotConnected
	}
	return nil
}

// Get retrieves the value for the given key from the node.
func (node *CacheNode) Get(key any) (any, error) {
	if err := node.ping(); err != nil {
		return nil, err
	}

	val, found := node.data.Load(key)
	if !found {
		err := fmt.Errorf("key: %v, %w", key, ErrKeyNotFound)
		log.Printf("[Node: %s] ❌ GET key: %v (not found)\n", node.Identifier, key)
		return nil, err
	}

	log.Printf("[Node: %s] ✅ GET key: %v → %v\n", node.Identifier, key, val)
	return val, nil
}

// Put stores a key-value pair in the node.
func (node *CacheNode) Put(key, val any) error {
	if err := node.ping(); err != nil {
		return err
	}
	node.data.Store(key, val)
	log.Printf("[Node: %s] ➕ PUT key: %v → %v\n", node.Identifier, key, val)
	return nil
}

// Delete removes the key from the node.
func (node *CacheNode) Delete(key any) error {
	if err := node.ping(); err != nil {
		return err
	}
	val, found := node.data.LoadAndDelete(key)
	if !found {
		log.Printf("[Node: %s] ❌ 🗑️ DELETE key: %v (not found)\n", node.Identifier, key)
	} else {
		log.Printf("[Node: %s] ✅ 🗑️ DELETE key: %v → %v\n", node.Identifier, key, val)
	}
	return nil
}
