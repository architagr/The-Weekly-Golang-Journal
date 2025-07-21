package ginhandler

import (
	"net/http"
	"vectory_clock/key-value-node/internal/controller"
	"vectory_clock/pkg/model"

	"github.com/gin-gonic/gin"
)

func InitRouters(ginEngine *gin.Engine, ctrl *controller.Store) {
	// Initialize the router with the gin engine
	// This function will set up all the routes and handlers
	// for the key-value store application.

	ginEngine.GET("/:key", func(c *gin.Context) {
		// Handler for GET request to retrieve a value by key
		key := c.Param("key")
		// Logic to retrieve the value for the key
		v, ok := ctrl.Get(key)
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "Key not found"})
			return
		}
		c.JSON(http.StatusOK, v)
	})
	ginEngine.PUT("/:key", func(c *gin.Context) {
		// Handler for PUT request to set a value for a key
		key := c.Param("key")
		var value *model.ValueWithClock
		if err := c.ShouldBindJSON(&value); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}
		value = ctrl.Set(key, value)
		c.JSON(http.StatusOK, value)
	})

	// Add more routes as needed
}
