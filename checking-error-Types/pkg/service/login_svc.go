package service

import (
	apperrors "checking-error-types/pkg/app-errors"
	"checking-error-types/pkg/dto"
	"checking-error-types/pkg/entities"
	"checking-error-types/pkg/persistence"
	"fmt"
	"log"
)

var LoginServiceObj LoginService

func init() {
	LoginServiceObj = LoginService{
		loginPersistenceObj: persistence.LoginPersistenceObj,
	}
}

type ILoginPersistence interface {
	AuthenticateUser(authReq *dto.AuthRequest) (*entities.User, error)
}

type LoginService struct {
	loginPersistenceObj ILoginPersistence
}

func (svc LoginService) AuthenticateUser(authReq *dto.AuthRequest) (*dto.AuthResponse, error) {
	userInfo, err := svc.loginPersistenceObj.AuthenticateUser(authReq)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("request (%+v), %w", authReq, err)
	}
	activeSessions := userInfo.Sessions.Filter(func(val *entities.Session) bool {
		return val.EndTime.IsZero()
	})
	if len(activeSessions) > 0 {
		return nil, apperrors.ActiveSessionError{}
	}
	return (&dto.AuthResponse{}).Init(userInfo), nil
}
