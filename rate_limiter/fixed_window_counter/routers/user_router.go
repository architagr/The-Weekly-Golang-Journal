package routers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func init() {
	UserRouterObj = UserRouter{}
}

var UserRouterObj UserRouter

type UserRouter struct{}

func (r *UserRouter) Get(c *gin.Context) {
	userID := c.Param("id")
	log.Printf("[handler] 📦 User fetch requested | ID: %s | IP: %s", userID, c.ClientIP())
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("User %s retrieved successfully", userID),
	})
}
