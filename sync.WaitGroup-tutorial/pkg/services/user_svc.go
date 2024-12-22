package service

import (
	"context"
	"fmt"
	"sync"
	"waitGroup-tutorial/pkg/dto"
	"waitGroup-tutorial/pkg/entities"
	"waitGroup-tutorial/pkg/persistence"
)

var UserServiceObj UserService

func init() {
	UserServiceObj = UserService{
		userPersistenceObj: &persistence.UserPersistenceObj,
	}
}

type IUserPersistence interface {
	Save(ctx context.Context, user *entities.User) (*entities.User, error)
}

type UserService struct {
	userPersistenceObj IUserPersistence
}

func (svc UserService) Signup(ctx context.Context, user *dto.User) (*dto.User, error) {
	fmt.Println("get in the svc")
	userInfo, err := svc.userPersistenceObj.Save(ctx, user.Map())
	if err != nil {
		return nil, fmt.Errorf("%w, user is not Saved", err)
	}
	// user is saved in DB
	wg := sync.WaitGroup{}
	wg.Add(1)
	go svc.sendWelcomeEmail(user, &wg)
	wg.Add(1)
	go svc.registerTags(user, &wg)
	wg.Wait()
	return (&dto.User{}).Init(userInfo), nil
}
func (svc UserService) sendWelcomeEmail(user *dto.User, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Sending info to a queue so that it is picked up by email sending microservice")
}

func (svc UserService) registerTags(user *dto.User, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Sending info to a queue so that it is picked up by notification service to register this user for tags")
}
