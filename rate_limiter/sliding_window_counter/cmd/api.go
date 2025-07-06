package main

import (
	"log"
	"sliding_window_counter/middlewares"
	"sliding_window_counter/routers"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("[server] Starting server on port 8081")
	engine := gin.New()

	engine.Use(middlewares.SlidingWindowCounterMiddleware(5, time.Minute, 6)) // precision: 6 sub-windows

	engine.GET("/v1/user/:id", routers.UserRouterObj.Get)
	if err := engine.Run(":8081"); err != nil {
		log.Fatalf("[server] Failed to start server: %v", err)
	}
}
