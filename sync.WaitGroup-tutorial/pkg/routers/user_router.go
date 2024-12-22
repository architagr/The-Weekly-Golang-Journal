package routers

import (
	"context"
	"net/http"
	"waitGroup-tutorial/pkg/dto"
	service "waitGroup-tutorial/pkg/services"

	"github.com/gin-gonic/gin"
)

type IUserService interface {
	Signup(ctx context.Context, user *dto.User) (*dto.User, error)
}

func init() {
	UserRouterObj = UserRouter{
		userServiceObj: service.UserServiceObj,
	}
}

var UserRouterObj UserRouter

type UserRouter struct {
	userServiceObj IUserService
}

func (r UserRouter) Signup(c *gin.Context) {
	var userReq dto.User
	if err := c.ShouldBindJSON(&userReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userInfo, err := r.userServiceObj.Signup(c.Request.Context(), &userReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, userInfo)
}
