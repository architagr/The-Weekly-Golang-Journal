package dto

import (
	"checking-error-types/pkg/entities"
	"time"
)

type User struct {
	Id          int
	Name        string
	CreatedDate time.Time
}

func (res *User) Init(userInfo *entities.User) *User {
	res = new(User)
	res.Id = userInfo.Id
	res.Name = userInfo.Name
	return res
}
