package ginhandler

import (
	"net/http"
	"vectory_clock/key-value-node/internal/controller"
	"vectory_clock/pkg/model"

	"github.com/gin-gonic/gin"
)

// InitRouters sets up all HTTP routes for this node.
func InitRouters(ginEngine *gin.Engine, ctrl *controller.Store) {
	// GET /:key - retrieve value by key
	ginEngine.GET("/:key", func(c *gin.Context) {
		key := c.Param("key")
		v, ok := ctrl.Get(key)
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "Key not found"})
			return
		}
		c.JSON(http.StatusOK, v)
	})

	// PUT /:key - set value for key with vector clock payload
	ginEngine.PUT("/:key", func(c *gin.Context) {
		key := c.Param("key")
		var value *model.ValueWithClock
		if err := c.ShouldBindJSON(&value); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}
		result := ctrl.Set(key, value)
		c.JSON(http.StatusOK, result)
	})
}
