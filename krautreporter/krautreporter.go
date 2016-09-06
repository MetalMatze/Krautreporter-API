package krautreporter

import (
	"github.com/go-kit/kit/log"
	"github.com/jinzhu/gorm"
	"github.com/metalmatze/krautreporter-api/krautreporter/interactor"
	"github.com/metalmatze/krautreporter-api/krautreporter/repository"
	gocache "github.com/patrickmn/go-cache"
)

// Krautreporter has all domain objects and dependencies
type Krautreporter struct {
	CrawlInteractor *interactor.CrawlInteractor
	HTTPInteractor  *interactor.HTTPInteractor
	Logger          log.Logger
}

// New returns a Krautreporter domain object
func New(logger log.Logger, db *gorm.DB, cache *gocache.Cache) *Krautreporter {

	authorRepository := repository.NewGormAuthorRepository(logger, db, cache)
	articleRepository := repository.NewGormArticleRepository(logger, db, cache)
	crawlRepository := repository.NewCrawlRepository(logger, db, cache)

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
		Logger: logger,
	}
}
