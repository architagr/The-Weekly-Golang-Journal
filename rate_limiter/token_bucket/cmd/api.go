package main

import (
	"log"
	"token_bucket/middlewares"
	"token_bucket/routers"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("[server] Starting server on port 8081")
	engine := gin.New()

	engine.Use(middlewares.RateLimiterMiddleware(4, 4))

	engine.GET("/v1/user/:id", routers.UserRouterObj.Get)
	if err := engine.Run(":8081"); err != nil {
		log.Fatalf("[server] Failed to start server: %v", err)
	}
}
