package ginhandler

import (
	"log"
	"net/http"
	"vectory_clock/key-value-store/internal/controller"
	"vectory_clock/key-value-store/internal/gateway"
	"vectory_clock/pkg/model"

	"github.com/gin-gonic/gin"
)

type clusterRouteHandler struct {
	ctrl *controller.Cluster
}

func NewClusterRouteHandler(ctrl *controller.Cluster) *clusterRouteHandler {
	return &clusterRouteHandler{ctrl: ctrl}
}

func (h *clusterRouteHandler) GetValue(c *gin.Context) {
	// Handler for GET request to retrieve a value by key
	key := c.Param("key")
	v, err := h.ctrl.Get(key)
	if err != nil {
		log.Printf("Error getting value for key %s: %v", key, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Key not found"})
		return
	}
	c.JSON(http.StatusOK, v)
}
func (h *clusterRouteHandler) SetValue(c *gin.Context) {
	// Handler for PUT request to set a value for a key
	key := c.Param("key")
	var value *model.ValueWithClock
	if err := c.ShouldBindJSON(&value); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	value, err := h.ctrl.Set(key, value)
	if err != nil {
		log.Printf("Error setting value for key %s: %v", key, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set value"})
		return
	}
	log.Printf("Successfully set value for key %s: %v", key, value)
	c.JSON(http.StatusOK, value)
}

func (h *clusterRouteHandler) RegisterNode(c *gin.Context) {
	// Handler for registering a new node
	var node model.Node
	if err := c.ShouldBindJSON(&node); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid node data"})
		return
	}
	if node.ID == "" || node.Address == "" || node.Port <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Node ID, address, and port must be provided"})
		return
	}
	log.Printf("Registering node: %s at %s:%d", node.ID, node.Address, node.Port)
	// Register the node with the controller
	gNode, err := gateway.NewNode(node.ID, node.Address, node.Port)
	if err != nil {
		log.Printf("Error creating node: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create node"})
		return
	}
	if err := h.ctrl.AddNode(gNode); err != nil {
		log.Printf("Error registering node %s: %v", node.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register node	"})
		return
	}
	log.Printf("Successfully registered node %s", node.ID)
	c.JSON(http.StatusOK, gin.H{"message": "Node registered successfully"})
}

func (h *clusterRouteHandler) DeregisterNode(c *gin.Context) {
	// Handler for unregistering a node
	var node model.Node
	if err := c.ShouldBindJSON(&node); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid node data"})
		return
	}
	if node.ID == "" || node.Address == "" || node.Port <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Node ID, address, and port must be provided"})
		return
	}
	log.Printf("Unregistering node: %s", node.ID)
	gNode, err := gateway.NewNode(node.ID, node.Address, node.Port)
	if err != nil {
		log.Printf("Error creating node: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create node"})
		return
	}
	// Unregister the node from the controller
	if err := h.ctrl.RemoveNode(gNode); err != nil {
		log.Printf("Error unregistering node %s: %v", node.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unregister node"})
		return
	}
	log.Printf("Successfully unregistered node %s", node.ID)
	c.JSON(http.StatusOK, gin.H{"message": "Node unregistered successfully"})
}

func InitRouters(ginEngine *gin.Engine, ctrl *controller.Cluster) {
	// Initialize the router with the gin engine
	// This function will set up all the routes and handlers
	// for the key-value store application.
	h := NewClusterRouteHandler(ctrl)
	{
		ginEngine.GET("/:key", h.GetValue)
		ginEngine.PUT("/:key", h.SetValue)
	}

	nodeRoutes := ginEngine.Group("/node")
	{
		nodeRoutes.POST("/register", h.RegisterNode)
		nodeRoutes.POST("/deregister", h.DeregisterNode)
	}
}
