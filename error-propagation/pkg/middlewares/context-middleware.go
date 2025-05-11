package middlewares

import (
	"net/http"
	"time"

	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
)

// this is imitating the API gateway
func ContextMiddleware(t time.Duration) gin.HandlerFunc {
	return timeout.New(
		timeout.WithTimeout(t),
		timeout.WithHandler(func(c *gin.Context) {
			c.Next()
		}),
		timeout.WithResponse(func(c *gin.Context) {
			response := gin.H{"error": "gateway time out"}
			c.JSON(http.StatusGatewayTimeout, response)
		}),
	)
}
