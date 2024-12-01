package service

import (
	"checking-error-types/pkg/dto"
	"checking-error-types/pkg/entities"
	"checking-error-types/pkg/persistence"
	"fmt"
)

var UserServiceObj UserService

func init() {
	UserServiceObj = UserService{
		userPersistenceObj: persistence.UserPersistenceObj,
	}
}

type IUserPersistence interface {
	Get(id int) (*entities.User, error)
}

type UserService struct {
	userPersistenceObj IUserPersistence
}

func (svc UserService) Get(id int) (*dto.User, error) {
	userInfo, err := svc.userPersistenceObj.Get(id)
	if err != nil {
		return nil, fmt.Errorf("%w, entity that is not found is user", err)
	}
	return (&dto.User{}).Init(userInfo), nil
}
