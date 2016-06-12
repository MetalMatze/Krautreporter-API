package repository

import (
	"time"

	"github.com/MetalMatze/Krautreporter-API/krautreporter/entity"
	"github.com/jinzhu/gorm"
)

type CrawlRepository struct {
	DB *gorm.DB
}

func (r CrawlRepository) FindOutdatedAuthors() ([]entity.Author, error) {
	var crawls []*entity.Crawl
	r.DB.Where("next < ?", time.Now()).Where("crawlable_type = ?", "authors").Order("next").Find(&crawls)

	var IDs []int
	for _, c := range crawls {
		IDs = append(IDs, c.CrawlableID)
	}

	var authors []entity.Author
	r.DB.Preload("Crawl").Where(IDs).Find(&authors)

	return authors, nil
}

func (r CrawlRepository) FindOutdatedArticles() ([]entity.Article, error) {
	var crawls []*entity.Crawl
	r.DB.Where("next < ?", time.Now()).Where("crawlable_type = ?", "articles").Order("next").Find(&crawls)

	var IDs []int
	for _, c := range crawls {
		IDs = append(IDs, c.CrawlableID)
	}

	var articles []entity.Article
	r.DB.Preload("Crawl").Where(IDs).Find(&articles)

	return articles, nil
}

// NextCrawls gets the next crawls limited by number
func (r CrawlRepository) NextCrawls(limit int) ([]*entity.Crawl, error) {
	var crawls []*entity.Crawl

	if result := r.DB.Limit(limit).Order("next").Find(&crawls); result.Error != nil {
		return nil, result.Error
	}

	return crawls, nil
}
