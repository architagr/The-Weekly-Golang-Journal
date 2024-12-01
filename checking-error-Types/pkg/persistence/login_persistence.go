package persistence

import (
	apperrors "checking-error-types/pkg/app-errors"
	"checking-error-types/pkg/dto"
	"checking-error-types/pkg/entities"
	"time"

	"github.com/architagr/golang_collections/list"
)

var LoginPersistenceObj LoginPersistence

func init() {
	LoginPersistenceObj = LoginPersistence{}
}

type LoginPersistence struct {
}

func (svc LoginPersistence) AuthenticateUser(authReq *dto.AuthRequest) (*entities.User, error) {
	if authReq.UserName == "user2" && authReq.Password == "password2" {
		// this user has an active session
		userInfo := new(entities.User)
		userInfo.Id = 2
		userInfo.Name = "user-2"
		userInfo.Sessions = list.InitArrayList[*entities.Session](
			&entities.Session{
				Id:        1,
				UserId:    userInfo.Id,
				StartTime: time.Now().Add(-1 * time.Hour),
			},
		)
		return userInfo, nil
	} else if authReq.UserName == "user3" && authReq.Password == "password3" {
		// this user does not has any active session
		userInfo := new(entities.User)
		userInfo.Id = 3
		userInfo.Name = "user-3"
		userInfo.Sessions = list.InitArrayList[*entities.Session](
			&entities.Session{
				Id:        1,
				UserId:    userInfo.Id,
				StartTime: time.Now().Add(-1 * time.Hour),
				EndTime:   time.Now().Add(-1 * time.Minute),
			},
		)
		return userInfo, nil
	}

	return nil, apperrors.CredentialError{}
}
