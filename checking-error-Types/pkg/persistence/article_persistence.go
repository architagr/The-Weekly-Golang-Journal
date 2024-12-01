package persistence

import (
	apperrors "checking-error-types/pkg/app-errors"
	"checking-error-types/pkg/entities"
	"time"
)

var ArticlePersistenceObj ArticlePersistence

func init() {
	ArticlePersistenceObj = ArticlePersistence{}
}

type ArticlePersistence struct {
}

func (p ArticlePersistence) Get(id int) (*entities.Article, error) {
	if id == 1 {
		return &entities.Article{
			Id:          1,
			Title:       "Article 1",
			Description: "Article 1 description",
			CreatedAt:   time.Now().Add(-10 * time.Hour),
			CreatedBy:   1,
		}, nil
	}
	return nil, apperrors.NotFound{Id: id}
}
