package cacheserver

import (
	cache_node "consistent_hashing/cache_node"
	hashring "consistent_hashing/redundant_hashring"
	"errors"
	"fmt"
	"hash"
	"log"
	"math/rand"
)

type CachingServer struct {
	nodes    []*cache_node.CacheNode
	hashRing *hashring.HashRing
}

// InitCachingServer creates a new caching cluster with N cache nodes,
// registers them with the consistent hash ring, and simulates random failure.
func InitCachingServer(hashFunc func() hash.Hash64, countNodes int) *CachingServer {
	if countNodes <= 0 {
		return nil
	}

	hashRing := hashring.InitHashRing(
		hashring.SetHashFunction(hashFunc),
		hashring.EnableVerboseLogs(true),
		hashring.SetVirtualNodes(3),
	)

	nodes := make([]*cache_node.CacheNode, countNodes)
	for i := 0; i < countNodes; i++ {
		identifier := fmt.Sprintf("%c_%d_node_%d", 'a'+rand.Intn(26), rand.Intn(1000), i)
		log.Printf("[Node: %s] ⚠️  Will simulate failure randomly in the next 0–10 seconds\n", identifier)

		node := cache_node.InitCacheNode(identifier)
		nodes[i] = node
		if err := hashRing.AddNode(node); err != nil {
			log.Printf("[Init] Failed to add node %s: %v", identifier, err)
		}
	}

	return &CachingServer{
		nodes:    nodes,
		hashRing: hashRing,
	}
}

// Put inserts a key-value pair into the appropriate cache node.
// If the node is disconnected, it retries after removing the node from the ring.
func (c *CachingServer) Put(key string, val any) error {
	node, err := c.putOnce(key, val)
	if err == nil {
		return nil
	}

	if errors.Is(err, cache_node.ErrNodeNotConnected) {
		log.Printf("⚠️  Node %s disconnected during PUT. Retrying...", node.GetIdentifier())
		c.removeNode(node)
		_, retryErr := c.putOnce(key, val)
		return retryErr
	}
	return err
}

func (c *CachingServer) putOnce(key string, val any) (*cache_node.CacheNode, error) {
	nodes, err := c.getNodes(key)
	if err != nil {
		return nil, err
	}
	for _, node := range nodes {
		err := node.Put(key, val)
		if err != nil {
			return node, err
		}
	}

	return nodes[0], nil
}

// Get retrieves the value associated with the key.
// If the node is disconnected, it retries after removing the node from the ring.
func (c *CachingServer) Get(key string) (any, error) {
	node, val, err := c.getOnce(key)
	if err == nil {
		return val, nil
	}

	if errors.Is(err, cache_node.ErrNodeNotConnected) {
		log.Printf("⚠️  Node %s disconnected during GET. Retrying...", node.GetIdentifier())
		c.removeNode(node)
		node, val, retryErr := c.getOnce(key)
		if retryErr != nil {
			return val, fmt.Errorf("%w [node: %s]", retryErr, node.GetIdentifier())
		}
		return val, retryErr
	}
	return "", err
}

func (c *CachingServer) getOnce(key string) (*cache_node.CacheNode, any, error) {
	node, err := c.getNode(key)
	if err != nil {
		return node, "", err
	}
	val, err := node.Get(key)
	return node, val, err
}

// Delete removes the key from the appropriate cache node.
func (c *CachingServer) Delete(key string) error {
	nodes, err := c.getNodes(key)
	if err != nil {
		return err
	}
	for _, node := range nodes {
		err := node.Delete(key)
		if err != nil {
			return err
		}
	}
	return nil
}

// getNode uses the hash ring to determine the correct node for a given key.
func (c *CachingServer) getNode(key string) (*cache_node.CacheNode, error) {
	node, err := c.hashRing.GetPrimaryNode(key)
	if err != nil {
		return nil, err
	}
	return node.(*cache_node.CacheNode), nil
}

// getNode uses the hash ring to determine the correct node for a given key.
func (c *CachingServer) getNodes(key string) ([]*cache_node.CacheNode, error) {
	nodes, err := c.hashRing.GetNodesForKey(key)
	if err != nil {
		return nil, err
	}
	result := make([]*cache_node.CacheNode, len(nodes))
	for i, n := range nodes {
		result[i] = n.(*cache_node.CacheNode)
	}
	return result, nil
}

// removeNode removes a disconnected node from both the slice and the hash ring.
func (c *CachingServer) removeNode(node *cache_node.CacheNode) {
	if node == nil {
		return
	}
	log.Printf("🧹 Removing node: %s", node.GetIdentifier())

	// Remove from slice
	index := -1
	for i, n := range c.nodes {
		if n.GetIdentifier() == node.GetIdentifier() {
			index = i
			break
		}
	}
	if index >= 0 {
		c.nodes = append(c.nodes[:index], c.nodes[index+1:]...)
	}

	// Remove from hash ring
	if err := c.hashRing.RemoveNode(node); err != nil {
		log.Printf("❌ Failed to remove node from ring: %s → %v", node.GetIdentifier(), err)
	}
}
