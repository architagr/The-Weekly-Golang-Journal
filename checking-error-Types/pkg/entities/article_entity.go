package entities

import "time"

type Article struct {
	Id          int
	Title       string
	Description string
	CreatedAt   time.Time
	CreatedBy   int
}
