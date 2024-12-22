package dto

import (
	"time"
	"waitGroup-tutorial/pkg/entities"
)

type User struct {
	Id               int       `json:"id" uri:"id" binding:"required"`
	Name             string    `json:"name"`
	Email            string    `json:"email"`
	NotificationTags []string  `json:"notificationTags"`
	CreatedDate      time.Time `json:"-"`
}

func (res *User) Init(userInfo *entities.User) *User {
	res = new(User)
	res.Id = userInfo.Id
	res.Name = userInfo.Name
	res.Email = userInfo.Email
	res.NotificationTags = userInfo.NotificationTags
	return res
}
func (res *User) Map() *entities.User {
	userInfo := new(entities.User)
	userInfo.Id = res.Id
	userInfo.Name = res.Name
	userInfo.Email = res.Email
	userInfo.NotificationTags = res.NotificationTags
	return userInfo
}
