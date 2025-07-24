package controller

import (
	"gossip-protocol/config"
	"gossip-protocol/model"
	"log"
	"maps"
	"math/rand"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	nodeHealth   = map[string]bool{}
	nodeHealthMu sync.RWMutex
	configLock   sync.RWMutex
)

func StartHTTPServer() {
	// gin.SetMode(gin.DebugMode)
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, GetNodeHealth())
	})

	r.POST("/gossip", func(c *gin.Context) {
		var msg model.GossipMessage
		if err := c.ShouldBindJSON(&msg); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		log.Printf("Received gossip from %s\n", msg.SenderID)
		UpdateNodeHealth(msg.NodeHealth)
		UpdateConfigIfNewer(msg.Config)
		nodeHealth[msg.SenderID.String()] = true
		log.Printf("Updated node health: %v\n", nodeHealth)
		c.JSON(http.StatusOK, gin.H{"status": "received"})
	})

	r.POST("/join", func(c *gin.Context) {
		var peer struct {
			URL string `json:"url"`
		}
		if err := c.BindJSON(&peer); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		configLock.Lock()
		log.Println("New peer joined:", peer.URL)
		config.Peers = append(config.Peers, peer.URL)
		configLock.Unlock()
		c.JSON(http.StatusOK, gin.H{"message": "Peer added"})
	})

	r.Run(":" + config.Port)
}

func GetNodeHealth() map[string]bool {
	nodeHealthMu.RLock()
	defer nodeHealthMu.RUnlock()
	copy := make(map[string]bool)

	i := 0
	for k, v := range nodeHealth {
		copy[k] = v
		i++
		if i >= config.CurrentConfig.BufferSize {
			break
		}
	}

	return copy
}

func UpdateNodeHealth(newHealth map[string]bool) {
	nodeHealthMu.Lock()
	defer nodeHealthMu.Unlock()
	maps.Copy(nodeHealth, newHealth)
}

// UpdateConfigIfNewer updates the current gossip configuration with the values from the provided cfg.
// This function acquires a lock to ensure thread-safe updates to the configuration.
// Note: Timestamp validation is currently skipped, so the update is unconditional.
func UpdateConfigIfNewer(cfg config.GossipConfig) {
	// Skipping timestamp validation for now
	configLock.Lock()
	defer configLock.Unlock()
	config.CurrentConfig.Fanout = cfg.Fanout
	config.CurrentConfig.Interval = cfg.Interval
	config.CurrentConfig.BufferSize = cfg.BufferSize
}

// GetRandomPeers returns a slice containing up to n randomly selected peer addresses
// from the current configuration's list of peers. If n is greater than the number
// of available peers, all peers are returned in random order. The function acquires
// a read lock on the configuration to ensure thread-safe access.
func GetRandomPeers(n int) []string {
	configLock.RLock()
	defer configLock.RUnlock()
	selected := make([]string, 0, n)
	perm := rand.Perm(len(config.Peers))
	for i := 0; i < n && i < len(config.Peers); i++ {
		selected = append(selected, config.Peers[perm[i]])
	}
	return selected
}
