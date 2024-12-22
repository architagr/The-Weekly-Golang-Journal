package main

import (
	"time"
	"waitGroup-tutorial/pkg/middlewares"
	"waitGroup-tutorial/pkg/routers"

	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.New()

	engine.Use(middlewares.ContextMiddleware(1 * time.Minute))
	engine.POST("/v1/user", routers.UserRouterObj.Signup)
	engine.Run(":8081")
}
