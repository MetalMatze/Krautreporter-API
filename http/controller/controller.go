package controller

import (
	"github.com/MetalMatze/Krautreporter-API/krautreporter/entity"
	"github.com/gollection/gollection/log"
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
	log        log.Logger
}

func New(log log.Logger, interactor HTTPInteractor) *Controller {
	return &Controller{interactor, log}
}
