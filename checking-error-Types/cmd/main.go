package main

import (
	"checking-error-types/pkg/routers"

	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()
	engine.POST("/v1/auth", routers.LoginRouterObj.AuthenticateV1)
	engine.POST("/v2/auth", routers.LoginRouterObj.AuthenticateV2)
	engine.Run(":8081")
}
