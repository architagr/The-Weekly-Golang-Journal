package main

import (
	"error-propagation/pkg/middlewares"
	"error-propagation/pkg/routers"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.New()
	engine.Use(middlewares.ContextMiddleware(10000 * time.Millisecond))

	engine.GET("/v1/user/:id", routers.UserRouterObj.Get)
	engine.Run(":8081")
}
