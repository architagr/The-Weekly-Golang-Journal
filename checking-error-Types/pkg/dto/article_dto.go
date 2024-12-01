package dto

import (
	"checking-error-types/pkg/entities"
	"time"
)

type Article struct {
	Id          int
	Title       string
	Description string
	CreatedAt   time.Time
	CreatedBy   int
}

func (res *Article) Init(articleInfo *entities.Article) *Article {
	res = new(Article)
	res.Id = articleInfo.Id
	res.Title = articleInfo.Title
	res.Description = articleInfo.Description
	res.CreatedAt = articleInfo.CreatedAt
	res.CreatedBy = articleInfo.CreatedBy
	return res
}
