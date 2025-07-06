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
	log.Printf("[handler] Fetching user with ID: %s", userID)
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("user %s retrieved", userID),
	})
}
