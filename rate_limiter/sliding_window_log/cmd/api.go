package main

import (
	"log"
	"sliding_window_log/middlewares"
	"sliding_window_log/routers"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("[server] Starting server on port 8081")
	engine := gin.New()

	engine.Use(middlewares.SlidingWindowLogMiddleware(2, 1*time.Minute))

	engine.GET("/v1/user/:id", routers.UserRouterObj.Get)
	if err := engine.Run(":8081"); err != nil {
		log.Fatalf("[server] Failed to start server: %v", err)
	}
}
