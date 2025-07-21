package main

import (
	"flag"
	"log"
	"vectory_clock/key-value-store/internal/controller"
	"vectory_clock/key-value-store/internal/handler/ginhandler"

	"github.com/gin-gonic/gin"
)

var clstr *controller.Cluster

var (
	readQuorum    = flag.Int("read-quorum", 2, "Number of nodes required to read a value (R)")
	writeQuorum   = flag.Int("write-quorum", 2, "Number of nodes required to write a value (W)")
	totalReplicas = flag.Int("total-replicas", 3, "Total number of replicas in the cluster (N)")
	virtualNodes  = flag.Int("virtual-nodes", 3, "Number of virtual nodes per physical node")
)

func init() {
	flag.Parse()
	log.Printf("[CONFIG] R=%d, W=%d, N=%d, virtual-nodes=%d", *readQuorum, *writeQuorum, *totalReplicas, *virtualNodes)
	c, err := controller.NewCluster(
		controller.WithReadQuorum(*readQuorum),
		controller.WithWriteQuorum(*writeQuorum),
		controller.WithTotalReplicas(*totalReplicas),
		controller.WithVirtualNodes(*virtualNodes),
	)
	if err != nil {
		log.Fatalf("Failed to initialize cluster controller: %v", err)
	}
	clstr = c
}

func main() {
	gin.SetMode(gin.DebugMode)
	engine := gin.Default()
	ginhandler.InitRouters(engine, clstr)
	log.Println("[INFO] Key-Value Store (API) running on :8080")
	if err := engine.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
