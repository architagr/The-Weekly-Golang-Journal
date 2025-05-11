package routers

import (
	"context"
	"error-propagation/pkg/dto"
	service "error-propagation/pkg/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IHttpResponseError interface {
	GetHttpStatusCode() int
	GetErrorResponse() *dto.ErrorResponse
}

type IUserService interface {
	Get(ctx context.Context, id int) (*dto.User, error)
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

func (r UserRouter) Get(c *gin.Context) {
	var userReq dto.User
	if err := c.ShouldBindUri(&userReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userInfo, err := r.userServiceObj.Get(c.Request.Context(), userReq.Id)
	if err != nil {
		if e, ok := err.(IHttpResponseError); ok {
			c.JSON(e.GetHttpStatusCode(), e.GetErrorResponse())
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, userInfo)
}
