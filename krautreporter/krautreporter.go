package krautreporter

import (
	"github.com/MetalMatze/Krautreporter-API/krautreporter/interactor"
	"github.com/MetalMatze/Krautreporter-API/krautreporter/repository"
	"github.com/gollection/gollection"
	"github.com/gollection/gollection/log"
)

// Krautreporter has all domain objects and dependencies
type Krautreporter struct {
	CrawlInteractor *interactor.CrawlInteractor
	HTTPInteractor  *interactor.HTTPInteractor
	Log             log.Logger
}

// New returns a Krautreporter domain object
func New(g *gollection.Gollection) *Krautreporter {

	authorRepository := repository.NewGormAuthorRepository(g.Cache, g.DB, g.Log)
	articleRepository := repository.NewGormArticleRepository(g.Cache, g.DB, g.Log)
	crawlRepository := repository.NewCrawlRepository(g.Cache, g.DB, g.Log)

	return &Krautreporter{
		CrawlInteractor: interactor.NewCrawlInteractor(
			authorRepository,
			articleRepository,
			crawlRepository,
		),
		HTTPInteractor: interactor.NewHTTPInteractor(
			authorRepository,
			articleRepository,
			crawlRepository,
		),
		Log: g.Log,
	}
}
