package service

import (
	"checking-error-types/pkg/dto"
	"checking-error-types/pkg/entities"
	"checking-error-types/pkg/persistence"
	"fmt"
)

var ArticleServiceObj ArticleService

func init() {
	ArticleServiceObj = ArticleService{
		articlePersistenceObj: persistence.ArticlePersistenceObj,
	}
}

type IArticlePersistence interface {
	Get(id int) (*entities.Article, error)
}

type ArticleService struct {
	articlePersistenceObj IArticlePersistence
}

func (svc ArticleService) Get(id int) (*dto.Article, error) {
	userInfo, err := svc.articlePersistenceObj.Get(id)
	if err != nil {
		return nil, fmt.Errorf("%w, entity that is not found is article", err)
	}
	return (&dto.Article{}).Init(userInfo), nil
}
