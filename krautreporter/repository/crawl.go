package repository

import (
	"time"

	"github.com/MetalMatze/Krautreporter-API/krautreporter/entity"
	"github.com/go-kit/kit/log"
	"github.com/jinzhu/gorm"
	gocache "github.com/patrickmn/go-cache"
)

type CrawlRepository struct {
	repository
}

func NewCrawlRepository(logger log.Logger, db *gorm.DB, cache *gocache.Cache) *CrawlRepository {
	return &CrawlRepository{repository: newRepository(logger, db, cache)}
}

func (r CrawlRepository) FindOutdatedAuthors() ([]entity.Author, error) {
	var crawls []*entity.Crawl
	r.db.Where("next < ?", time.Now()).Where("crawlable_type = ?", "authors").Order("next").Find(&crawls)

	var IDs []int
	for _, c := range crawls {
		IDs = append(IDs, c.CrawlableID)
	}

	var authors []entity.Author
	r.db.Preload("Crawl").Where(IDs).Find(&authors)

	return authors, nil
}

func (r CrawlRepository) FindOutdatedArticles() ([]entity.Article, error) {
	var crawls []*entity.Crawl
	r.db.Where("next < ?", time.Now()).Where("crawlable_type = ?", "articles").Order("next").Find(&crawls)

	var IDs []int
	for _, c := range crawls {
		IDs = append(IDs, c.CrawlableID)
	}

	var articles []entity.Article
	r.db.Preload("Crawl").Where(IDs).Find(&articles)

	return articles, nil
}

// NextCrawls gets the next crawls limited by number
func (r CrawlRepository) NextCrawls(limit int) ([]*entity.Crawl, error) {
	var crawls []*entity.Crawl

	if result := r.db.Limit(limit).Order("next").Find(&crawls); result.Error != nil {
		return nil, result.Error
	}

	return crawls, nil
}
