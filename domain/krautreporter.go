package domain

import (
	"github.com/MetalMatze/Krautreporter-API/domain/interactor"
	"github.com/MetalMatze/Krautreporter-API/domain/repository"
	"github.com/MetalMatze/gollection"
)

type Krautreporter struct {
	AuthorInteractor  interactor.AuthorInteractor
	ArticleInteractor interactor.ArticleInteractor
}

func NewKrautreporter(g *gollection.Gollection) *Krautreporter {
	return &Krautreporter{
		AuthorInteractor: interactor.AuthorInteractor{
			AuthorRepository: repository.GormAuthorRepository{
				Cache: g.Cache,
				DB:    g.DB,
			},
		},
		ArticleInteractor: interactor.ArticleInteractor{
			ArticleRepository: repository.GormArticleRepository{
				DB: g.DB,
			},
		},
	}
}
