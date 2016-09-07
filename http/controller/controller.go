package controller

import (
	"github.com/go-kit/kit/log"
	"github.com/metalmatze/krautreporter-api/entity"
)

type HTTPInteractor interface {
	AllAuthors() ([]*entity.Author, error)
	AuthorByID(id int) (*entity.Author, error)

	ArticlesOlderThan(id int, number int) ([]*entity.Article, error)
	ArticleByID(id int) (*entity.Article, error)

	NextCrawls(limit int) ([]*entity.Crawl, error)
}

type Controller struct {
	interactor HTTPInteractor
	logger     log.Logger
}

func New(logger log.Logger, interactor HTTPInteractor) *Controller {
	return &Controller{interactor, logger}
}
