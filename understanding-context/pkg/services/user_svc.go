package service

import (
	"context"
	"fmt"
	"understanding-context/pkg/dto"
	"understanding-context/pkg/entities"
	"understanding-context/pkg/persistence"
)

var UserServiceObj UserService

func init() {
	UserServiceObj = UserService{
		userPersistenceObj: &persistence.UserPersistenceObj,
	}
}

type IUserPersistence interface {
	Get(ctx context.Context, id int) (*entities.User, error)
}

type UserService struct {
	userPersistenceObj IUserPersistence
}

func (svc UserService) Get(ctx context.Context, id int) (*dto.User, error) {
	fmt.Println("get in the svc")
	userInfo, err := svc.userPersistenceObj.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w, user is not found", err)
	}
	return (&dto.User{}).Init(userInfo), nil
}
