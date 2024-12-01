package routers

import (
	apperrors "checking-error-types/pkg/app-errors"
	"checking-error-types/pkg/dto"
	"checking-error-types/pkg/service"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ILoginService interface {
	AuthenticateUser(authReq *dto.AuthRequest) (*dto.AuthResponse, error)
}

func init() {
	LoginRouterObj = LoginRouter{
		loginServiceObj: service.LoginServiceObj,
	}
}

var LoginRouterObj LoginRouter

type LoginRouter struct {
	loginServiceObj ILoginService
}

func (r LoginRouter) AuthenticateV1(c *gin.Context) {
	var authReq dto.AuthRequest
	if err := c.ShouldBindJSON(&authReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userInfo, err := r.loginServiceObj.AuthenticateUser(&authReq)
	if err != nil {
		switch err := err.(type) {
		case apperrors.CredentialError:
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		case apperrors.ActiveSessionError:
			c.JSON(http.StatusPreconditionFailed, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, userInfo)
}

func (r LoginRouter) AuthenticateV2(c *gin.Context) {
	var authReq dto.AuthRequest
	if err := c.ShouldBindJSON(&authReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userInfo, err := r.loginServiceObj.AuthenticateUser(&authReq)
	if err != nil {
		var credErr apperrors.CredentialError
		var sessionErr apperrors.ActiveSessionError
		if errors.As(err, &apperrors.CredentialError{}) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": credErr.Error()})
		} else if errors.As(err, &apperrors.ActiveSessionError{}) {
			c.JSON(http.StatusPreconditionFailed, gin.H{"error": sessionErr.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, userInfo)
}
