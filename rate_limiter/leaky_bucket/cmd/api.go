package main

import (
	"leaky_bucket/middlewares"
	"leaky_bucket/routers"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("[server] Starting server on port 8081")
	engine := gin.New()

	engine.Use(middlewares.LeakyBucketMiddleware(1, 1))

	engine.GET("/v1/user/:id", routers.UserRouterObj.Get)
	if err := engine.Run(":8081"); err != nil {
		log.Fatalf("[server] Failed to start server: %v", err)
	}
}
