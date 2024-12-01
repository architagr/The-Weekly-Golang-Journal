package persistence

import (
	apperrors "checking-error-types/pkg/app-errors"
	"checking-error-types/pkg/entities"
	"time"
)

var UserPersistenceObj UserPersistence

func init() {
	UserPersistenceObj = UserPersistence{}
}

type UserPersistence struct {
}

func (p UserPersistence) Get(id int) (*entities.User, error) {
	if id == 1 {
		return &entities.User{
			Id:          1,
			Name:        "user-1",
			CreatedDate: time.Now().Add(-10 * time.Hour),
		}, nil
	}
	return nil, apperrors.NotFound{Id: id}
}
