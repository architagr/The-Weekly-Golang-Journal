package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"vectory_clock/key-value-node/internal/config"
	"vectory_clock/key-value-node/internal/controller"
	"vectory_clock/key-value-node/internal/handler/ginhandler"
	globalModel "vectory_clock/pkg/model"

	"github.com/gin-gonic/gin"
)

var (
	nodeId               = flag.String("node-id", "", "Unique identifier for this node")
	address              = flag.String("address", "localhost", "Address of the node")
	port                 = flag.Int("port", 8080, "Port of the node")
	keyValueStoreAddress = flag.String("key-value-store-address", "localhost", "Address of the key-value store")
	keyValueStorePort    = flag.Int("key-value-store-port", 8080, "Port of the key-value store")
)

func init() {
	flag.Parse()
	if *nodeId == "" || *address == "" || port == nil || *port <= 0 {
		log.Fatalf("[FATAL] node-id, address, port are required to start the node.")
	}
	if *keyValueStoreAddress == "" || keyValueStorePort == nil || *keyValueStorePort <= 0 {
		log.Fatalf("[FATAL] key-value-store-address and key-value-store-port are required.")
	}
	config.NodeId = *nodeId
	registerWithKeyValueStore()
}

func registerWithKeyValueStore() {
	node := globalModel.Node{
		ID:      *nodeId,
		Address: *address,
		Port:    *port,
	}
	body, err := json.Marshal(node)
	if err != nil {
		log.Fatalf("[FATAL] Unable to serialize node info for registration: %v", err)
	}
	url := fmt.Sprintf("http://%s:%d/node/register", *keyValueStoreAddress, *keyValueStorePort)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		log.Fatalf("[FATAL] Unable to construct register request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("[FATAL] Unable to register node with key-value-store: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("[FATAL] Node registration failed: HTTP %d", resp.StatusCode)
	}
	log.Printf("[INFO] Node %s registered with key-value-store at %s:%d", node.ID, *keyValueStoreAddress, *keyValueStorePort)
}

func deregisterFromKeyValueStore() {
	node := globalModel.Node{
		ID:      *nodeId,
		Address: *address,
		Port:    *port,
	}
	body, err := json.Marshal(node)
	if err != nil {
		log.Printf("[WARN] Unable to serialize node info for deregistration: %v", err)
		return
	}
	url := fmt.Sprintf("http://%s:%d/node/deregister", *keyValueStoreAddress, *keyValueStorePort)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("[WARN] Unable to construct deregister request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("[WARN] Unable to deregister from key-value-store: %v", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("[WARN] Node deregistration failed: HTTP %d", resp.StatusCode)
	}
	log.Printf("[INFO] Node %s deregistered from key-value-store", node.ID)
}

func main() {
	defer deregisterFromKeyValueStore()
	gin.SetMode(gin.ReleaseMode) // Use ReleaseMode for production, DebugMode for verbose logs

	log.Printf("[INFO] Starting Key-Value Node (ID=%s, Address=%s, Port=%d)", *nodeId, *address, *port)

	ctrl := controller.NewStore()
	router := gin.Default()
	ginhandler.InitRouters(router, ctrl)

	serverAddr := fmt.Sprintf("%s:%d", *address, *port)
	log.Printf("[INFO] Listening for requests at %s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("[FATAL] Failed to start Gin server: %v", err)
	}
}
