package entities

import (
	"time"
)

type User struct {
	Id          int
	Name        string
	CreatedDate time.Time
}
