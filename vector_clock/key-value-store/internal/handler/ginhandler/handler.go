package ginhandler

import (
	"log"
	"net/http"
	"vectory_clock/key-value-store/internal/controller"
	"vectory_clock/key-value-store/internal/gateway"
	"vectory_clock/pkg/model"

	"github.com/gin-gonic/gin"
)

// clusterRouteHandler acts as the glue for all HTTP cluster operations.
type clusterRouteHandler struct {
	ctrl *controller.Cluster
}

func NewClusterRouteHandler(ctrl *controller.Cluster) *clusterRouteHandler {
	return &clusterRouteHandler{ctrl: ctrl}
}

// GET /:key
func (h *clusterRouteHandler) GetValue(c *gin.Context) {
	key := c.Param("key")
	v, err := h.ctrl.Get(key)
	if err != nil {
		log.Printf("[ERROR] GET key=%s: %v", key, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Key not found"})
		return
	}
	c.JSON(http.StatusOK, v)
}

// PUT /:key
func (h *clusterRouteHandler) SetValue(c *gin.Context) {
	key := c.Param("key")
	var value *model.ValueWithClock
	if err := c.ShouldBindJSON(&value); err != nil {
		log.Printf("[WARN] PUT invalid body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	result, err := h.ctrl.Set(key, value)
	if err != nil {
		log.Printf("[ERROR] PUT key=%s failed: %v", key, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set value"})
		return
	}
	log.Printf("[INFO] PUT key=%s: value=%v clock=%v", key, result.Value, result.Clock)
	c.JSON(http.StatusOK, result)
}

// POST /node/register
func (h *clusterRouteHandler) RegisterNode(c *gin.Context) {
	var node model.Node
	if err := c.ShouldBindJSON(&node); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid node data"})
		return
	}
	if node.ID == "" || node.Address == "" || node.Port <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Node ID, address, and port must be provided"})
		return
	}
	log.Printf("[INFO] Registering node: %s at %s:%d", node.ID, node.Address, node.Port)
	gNode, err := gateway.NewNode(node.ID, node.Address, node.Port)
	if err != nil {
		log.Printf("[ERROR] create node: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create node"})
		return
	}
	if err := h.ctrl.AddNode(gNode); err != nil {
		log.Printf("[ERROR] registering node %s: %v", node.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register node"})
		return
	}
	log.Printf("[SUCCESS] Registered node %s", node.ID)
	c.JSON(http.StatusOK, gin.H{"message": "Node registered successfully"})
}

// POST /node/deregister
func (h *clusterRouteHandler) DeregisterNode(c *gin.Context) {
	var node model.Node
	if err := c.ShouldBindJSON(&node); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid node data"})
		return
	}
	if node.ID == "" || node.Address == "" || node.Port <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Node ID, address, and port must be provided"})
		return
	}
	log.Printf("[INFO] Unregistering node: %s", node.ID)
	gNode, err := gateway.NewNode(node.ID, node.Address, node.Port)
	if err != nil {
		log.Printf("[ERROR] create node: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create node"})
		return
	}
	if err := h.ctrl.RemoveNode(gNode); err != nil {
		log.Printf("[ERROR] unregistering node %s: %v", node.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unregister node"})
		return
	}
	log.Printf("[SUCCESS] Unregistered node %s", node.ID)
	c.JSON(http.StatusOK, gin.H{"message": "Node unregistered successfully"})
}

func InitRouters(ginEngine *gin.Engine, ctrl *controller.Cluster) {
	h := NewClusterRouteHandler(ctrl)
	ginEngine.GET("/:key", h.GetValue)
	ginEngine.PUT("/:key", h.SetValue)
	nodeRoutes := ginEngine.Group("/node")
	nodeRoutes.POST("/register", h.RegisterNode)
	nodeRoutes.POST("/deregister", h.DeregisterNode)
}
