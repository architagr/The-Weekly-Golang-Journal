package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"gossip-protocol/config"
	"gossip-protocol/controller"
	"gossip-protocol/spread"
	"gossip-protocol/strategy"
	"log"
	"net/http"
	"strings"
	"time"
)

var (
	port  = flag.String("port", "8080", "Port to run the gossip protocol server on")
	peers = flag.String("peers", "", "Comma-separated list of peer addresses to connect to")
)

func init() {
	// Parse CLI flags early
	flag.Parse()

	if len(*port) == 0 {
		log.Fatal("❌ Port cannot be empty")
	}

	// Initialize peer list if provided
	if len(*peers) > 0 {
		config.Peers = strings.Split(*peers, ",")
		log.Printf("✅ Initial peers loaded: %v", config.Peers)
	}
}

func main() {
	// Load default gossip configuration
	cfg := config.DefaultConfig()
	config.Port = *port

	// Select gossip and spread strategies
	gossipStrategy := strategy.GetGossipStrategy(config.Strategy)
	spreadStrategy := spread.GetSpreadStrategy(config.SpreadMethod)

	// Start HTTP server in background
	go controller.StartHTTPServer()
	log.Printf("🚀 Gossip node %v started on port %s (Strategy: %s, Spread: %s)\n",
		config.SelfID, config.Port, config.Strategy, config.SpreadMethod)

	// Announce to peers that we are joining
	for _, peer := range config.Peers {
		go func(url string) {
			joinPayload := map[string]string{"url": "http://localhost:" + config.Port}
			data, _ := json.Marshal(joinPayload)

			resp, err := http.Post(url+"/join", "application/json", bytes.NewBuffer(data))
			if err != nil {
				log.Printf("⚠️ Failed to join peer %s: %v", url, err)
				return
			}
			resp.Body.Close()
			log.Printf("✅ Joined peer: %s", url)
		}(peer)
	}

	// Periodically gossip
	ticker := time.NewTicker(cfg.Interval)
	defer ticker.Stop()

	for range ticker.C {
		state := controller.GetNodeHealth()

		msg := gossipStrategy.GenerateMessage(state)
		peerList := controller.GetRandomPeers(cfg.Fanout)

		if len(peerList) > 0 {
			log.Printf("📢 Gossiping to %d peers. NodeHealth: %v", len(peerList), state)
			spreadStrategy.Spread(msg, peerList)
		} else {
			log.Println("⚠️ No peers available to spread gossip.")
		}
	}
}
