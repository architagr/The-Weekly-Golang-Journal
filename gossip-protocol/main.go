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
	flag.Parse()
	if len(*port) == 0 {
		panic("Port cannot be empty")
	}
	if len(*peers) > 0 {
		config.Peers = strings.Split(*peers, ",")
	}

}
func main() {

	cfg := config.DefaultConfig()

	config.Port = *port

	strategy := strategy.GetGossipStrategy(config.Strategy)
	spreadMethod := spread.GetSpreadStrategy(config.SpreadMethod)

	go controller.StartHTTPServer()
	log.Printf("Gossip protocol serverID %v started on port: %d\n", config.SelfID, config.Port)

	for _, peer := range config.Peers {
		go func(url string) {
			payload, _ := json.Marshal(map[string]string{
				"url": "http://localhost:" + config.Port,
			})
			http.Post(url, "application/json", bytes.NewBuffer(payload))
		}(peer + "/join")
	}

	ticker := time.NewTicker(cfg.Interval)
	defer ticker.Stop()

	for range ticker.C {
		state := controller.GetNodeHealth()
		msg := strategy.GenerateMessage(state)
		peers := controller.GetRandomPeers(cfg.Fanout)
		if len(peers) >= 0 {
			spreadMethod.Spread(msg, peers)
		} else {
			log.Println("No peers available to spread gossip")
		}
	}
}
