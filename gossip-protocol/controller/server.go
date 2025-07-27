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
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, GetNodeHealth())
	})

	// Gossip message receiver
	r.POST("/gossip", func(c *gin.Context) {
		var msg model.GossipMessage
		if err := c.ShouldBindJSON(&msg); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		log.Printf("[RECV] Gossip received from %s | Health: %v", msg.SenderID, msg.NodeHealth)

		UpdateNodeHealth(msg.NodeHealth)
		UpdateConfigIfNewer(msg.Config)
		nodeHealth[msg.SenderID.String()] = true
		log.Printf("[STATE] Node health updated: %v", nodeHealth)

		c.JSON(http.StatusOK, gin.H{"status": "received"})
	})

	// Peer join endpoint
	r.POST("/join", func(c *gin.Context) {
		var peer struct {
			URL string `json:"url"`
		}
		if err := c.BindJSON(&peer); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		configLock.Lock()
		config.Peers = append(config.Peers, peer.URL)
		configLock.Unlock()

		log.Printf("[JOIN] Peer added: %s | Current peers: %v", peer.URL, config.Peers)
		c.JSON(http.StatusOK, gin.H{"message": "Peer added"})
	})

	log.Printf("[HTTP] Listening on :%s", config.Port)
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

func UpdateConfigIfNewer(cfg config.GossipConfig) {
	configLock.Lock()
	defer configLock.Unlock()

	config.CurrentConfig.Fanout = cfg.Fanout
	config.CurrentConfig.Interval = cfg.Interval
	config.CurrentConfig.BufferSize = cfg.BufferSize

	log.Printf("[CONFIG] Updated local config: %+v", config.CurrentConfig)
}

func GetRandomPeers(n int) []string {
	configLock.RLock()
	defer configLock.RUnlock()

	if len(config.Peers) == 0 {
		return []string{}
	}

	selected := make([]string, 0, n)
	perm := rand.Perm(len(config.Peers))
	for i := 0; i < n && i < len(config.Peers); i++ {
		selected = append(selected, config.Peers[perm[i]])
	}
	return selected
}
