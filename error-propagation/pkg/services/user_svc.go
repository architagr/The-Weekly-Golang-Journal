package service

import (
	"context"
	"error-propagation/pkg/dto"
	"error-propagation/pkg/entities"
	"error-propagation/pkg/persistence"
	"errors"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
)

type ILogMessageError interface {
	error
	LogMessage() string
}

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
		if e, ok := err.(ILogMessageError); ok {
			return nil, svc.WrapError(e)
		}
		logBug(err)
		return nil, fmt.Errorf("can not retrive user: %d", id)
	}
	return (&dto.User{}).Init(userInfo), nil
}

// handle all types of errors that we expect the persistence to return
// or we can have a builder to build the ServiceError, and also follow
// The open and closed principle, and single responsibility principal
func (svc UserService) WrapError(err ILogMessageError) *ServiceError {
	var objNotFound *persistence.ObjectNotFoundError
	var result *ServiceError
	if errors.As(err, &objNotFound) {
		result = &ServiceError{
			InnerError:     err,
			HttpStatusCode: http.StatusNotFound,
			HttpResponse: &dto.ErrorResponse{
				ErrorCode:        "usr-404",
				ErrorDescription: err.Error(),
			},
		}
	}
	return result
}
func logBug(err error) {
	// log a bug and page enginnering team about the error
	log.Println("logged a bug for error", err, string(debug.Stack()))
}
