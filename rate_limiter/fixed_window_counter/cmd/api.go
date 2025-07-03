package main

import (
	"fixed_window_counter/middlewares"
	"fixed_window_counter/routers"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("[server] 🚀 Starting server on port 8081")

	engine := gin.New()
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())
	engine.Use(middlewares.FixedWindowCounterMiddleware(5, time.Second))

	engine.GET("/v1/user/:id", routers.UserRouterObj.Get)

	if err := engine.Run(":8081"); err != nil {
		log.Fatalf("[server] ❌ Failed to start server: %v", err)
	}
}
