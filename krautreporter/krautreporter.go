package krautreporter

import (
	"github.com/MetalMatze/Krautreporter-API/krautreporter/interactor"
	"github.com/MetalMatze/Krautreporter-API/krautreporter/repository"
	"github.com/gollection/gollection"
	"github.com/gollection/gollection/log"
)

// Krautreporter has all domain objects and dependencies
type Krautreporter struct {
	AuthorInteractor  interactor.AuthorInteractor
	ArticleInteractor interactor.ArticleInteractor
	CrawlInteractor   interactor.CrawlInteractor
	Log               log.Logger
}

// New returns a Krautreporter domain object
func New(g *gollection.Gollection) *Krautreporter {
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
				DB:  g.DB,
				Log: g.Log,
			},
		},
		CrawlInteractor: interactor.CrawlInteractor{
			AuthorRepository: authorRepository,
			CrawlRepository:  repository.CrawlRepository{DB: g.DB},
		},
		Log: g.Log,
	}
}
