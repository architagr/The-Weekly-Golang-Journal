package entities

import (
	"time"
)

type User struct {
	Id               int
	Name             string
	Email            string
	NotificationTags []string
	CreatedDate      time.Time
}
