package domain

import (
	"github.com/MetalMatze/Krautreporter-API/domain/interactor"
	"github.com/MetalMatze/Krautreporter-API/domain/repository"
	"github.com/gollection/gollection"
	"github.com/gollection/gollection/log"
)

type Krautreporter struct {
	AuthorInteractor  interactor.AuthorInteractor
	ArticleInteractor interactor.ArticleInteractor
	CrawlInteractor   interactor.CrawlInteractor
	Log               log.Logger
}

func NewKrautreporter(g *gollection.Gollection) *Krautreporter {
	authorRepository := repository.GormAuthorRepository{
		Cache: g.Cache,
		DB:    g.DB,
	}

	return &Krautreporter{
		AuthorInteractor: interactor.AuthorInteractor{
			AuthorRepository: authorRepository,
		},
		ArticleInteractor: interactor.ArticleInteractor{
			ArticleRepository: repository.GormArticleRepository{
				DB: g.DB,
			},
		},
		CrawlInteractor: interactor.CrawlInteractor{
			AuthorRepository: authorRepository,
			CrawlRepository:  repository.CrawlRepository{DB: g.DB},
		},
		Log: g.Log,
	}
}
