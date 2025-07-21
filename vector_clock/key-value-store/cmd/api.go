package main

import (
	"flag"
	"vectory_clock/key-value-store/internal/controller"
	"vectory_clock/key-value-store/internal/handler/ginhandler"

	"github.com/gin-gonic/gin"
)

var clstr *controller.Cluster

var (
	readQuorm     = flag.Int("read-quorum", 2, "Number of nodes required to read a value")
	writeQuorm    = flag.Int("write-quorum", 2, "Number of nodes required to write a value")
	totalReplicas = flag.Int("total-replicas", 3, "Total number of replicas in the cluster")
	virtualNodes  = flag.Int("virtual-nodes", 3, "Number of virtual nodes per physical node")
)

func init() {
	flag.Parse()

	// Initialize the cluster controller with the provided parameters
	c, err := controller.NewCluster(
		controller.WithReadQuorum(*readQuorm),
		controller.WithWriteQuorum(*writeQuorm),
		controller.WithTotalReplicas(*totalReplicas),
		controller.WithVirtualNodes(*virtualNodes),
	)
	if err != nil {
		panic("Failed to initialize cluster controller: " + err.Error())
	}
	clstr = c
}

func main() {
	gin.SetMode(gin.DebugMode)

	engine := gin.Default()
	// Start the HTTP server and initialize the routes
	ginhandler.InitRouters(engine, clstr)
	if err := engine.Run(":8080"); err != nil {
		panic("Failed to start server: " + err.Error())
	}
}
