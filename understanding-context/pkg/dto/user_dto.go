package dto

import (
	"time"
	"understanding-context/pkg/entities"
)

type User struct {
	Id          int       `json:"id" uri:"id" binding:"required"`
	Name        string    `json:"name"`
	CreatedDate time.Time `json:"-"`
}

func (res *User) Init(userInfo *entities.User) *User {
	res = new(User)
	res.Id = userInfo.Id
	res.Name = userInfo.Name
	return res
}
