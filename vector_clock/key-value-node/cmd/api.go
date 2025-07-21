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
	if *nodeId == "" {
		panic("node-id must be provided")
	}
	if *address == "" {
		panic("address must be provided")
	}
	if port == nil || *port <= 0 {
		panic("port must be provided")
	}

	if *keyValueStoreAddress == "" {
		panic("key-value-store-address must be provided")
	}

	if keyValueStorePort == nil || *keyValueStorePort <= 0 {
		panic("key-value-store-port must be provided")
	}

	// Set the global NodeId variable
	config.NodeId = *nodeId
	registerWithKeyValueStore()
}

func registerWithKeyValueStore() {
	// This function would typically register the node with the key-value store
	// For example, it could send a registration request to the key-value store's API
	// using the provided address and port.
	// This is a placeholder function, you would implement the actual registration logic here.
	node := globalModel.Node{
		ID:      *nodeId,
		Address: *address,
		Port:    *port,
	}
	body, err := json.Marshal(node)
	if err != nil {
		panic(fmt.Errorf("unable to register node: %w", err))
	}
	url := fmt.Sprintf("http://%s:%d/node/register", *keyValueStoreAddress, *keyValueStorePort)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		panic(fmt.Errorf("unable to register node: %w", err))
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(fmt.Errorf("unable to register node: %w", err))
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		panic(fmt.Errorf("unable to register node: %w", err))
	} else if resp.StatusCode != http.StatusOK {
		panic(fmt.Errorf("unable to register node: %w", err))
	}
}
func deregisterFromKeyValueStore() {
	// This function would typically deregister the node from the key-value store
	// For example, it could send a deregistration request to the key-value store's API
	// using the provided address and port.
	// This is a placeholder function, you would implement the actual deregistration logic here.
	node := globalModel.Node{
		ID:      *nodeId,
		Address: *address,
		Port:    *port,
	}
	body, err := json.Marshal(node)
	if err != nil {
		log.Panicln(fmt.Errorf("unable to deregister node: %w", err))
	}
	url := fmt.Sprintf("http://%s:%d/node/deregister", *keyValueStoreAddress, *keyValueStorePort)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		log.Panicln(fmt.Errorf("unable to deregister node: %w", err))
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Panicln(fmt.Errorf("unable to deregister node: %w", err))
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		log.Panicln(fmt.Errorf("unable to deregister node: %w", err))
	} else if resp.StatusCode != http.StatusOK {
		log.Panicln(fmt.Errorf("unable to deregister node: %w", err))
	}
}

func main() {
	// Start the key-value node server
	// This is where you would typically initialize your server and start listening for requests
	// For example:
	// server := NewServer(*address, *port)
	// if err := server.Start(); err != nil {
	//     log.Fatalf("Failed to start server: %v", err)
	// }
	defer deregisterFromKeyValueStore()
	gin.SetMode(gin.DebugMode)

	log.Printf("Key-Value Node started with ID: %s, Address: %s, Port: %d", *nodeId, *address, *port)
	// Here you would typically block the main goroutine to keep the server running
	ctrl := controller.NewStore()
	router := gin.Default()
	ginhandler.InitRouters(router, ctrl)
	if err := router.Run(fmt.Sprintf("%s:%d", *address, *port)); err != nil {
		log.Fatalf("Failed to start Gin server: %v", err)
	}
	// The server will now listen for incoming requests on the specified address and port
	// You can add more routes and handlers as needed
	log.Println("Server is running...")
}
