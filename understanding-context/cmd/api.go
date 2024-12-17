package main

import (
	"time"
	"understanding-context/pkg/middlewares"
	"understanding-context/pkg/routers"

	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.New()

	engine.Use(middlewares.ContextMiddleware(5000 * time.Millisecond))
	engine.GET("/v1/user/:id", routers.UserRouterObj.Get)
	engine.Run(":8081")
}
