package repository

import (
	"errors"
	"strings"
	"time"

	"github.com/MetalMatze/Krautreporter-API/krautreporter/entity"
	"github.com/go-kit/kit/log"
	"github.com/jinzhu/gorm"
	gocache "github.com/patrickmn/go-cache"
)

const MaxArticleID int = 1234567890

var ErrArticleNotFound = errors.New("Article not found")

type GormArticleRepository struct {
	repository
}

func NewGormArticleRepository(logger log.Logger, db *gorm.DB, cache *gocache.Cache) *GormArticleRepository {
	return &GormArticleRepository{repository: newRepository(logger, db, cache)}
}

func (r GormArticleRepository) FindOlderThan(id int, number int) ([]*entity.Article, error) {
	ordering := MaxArticleID
	if id != MaxArticleID {
		a, err := r.FindByID(id)
		if err != nil {
			return nil, err
		}

		ordering = a.Ordering
	}

	var articles []*entity.Article
	if result := r.db.
		Preload("Images").
		Where("ordering < ?", ordering).
		Limit(number).
		Order("ordering desc").
		Find(&articles); result.Error != nil {
		return nil, result.Error
	}

	return articles, nil
}

func (r GormArticleRepository) FindByID(id int) (*entity.Article, error) {
	var a entity.Article
	r.db.Preload("Images").Preload("Crawl").First(&a, "id = ?", id)

	if a.ID == 0 {
		return nil, ErrArticleNotFound
	}

	return &a, nil
}

func (r GormArticleRepository) SaveAll(articles []*entity.Article) error {
	tx := r.db.Begin()
	for i, a := range articles {
		article := entity.Article{ID: a.ID}
		tx.Preload("Images").Preload("Crawl").FirstOrCreate(&article)

		article.Ordering = len(articles) - 1 - i
		article.Title = a.Title
		article.URL = a.URL
		article.Preview = a.Preview

		author := entity.Author{}
		tx.First(&author, "name = ?", strings.TrimSpace(a.Author.Name))
		if author.ID == 0 {
			r.logger.Log("msg", "Can't find author for article ", "author", a.Author.Name, "article", a.URL)
			continue
		}
		article.AuthorID = author.ID

		for _, i := range a.Images {
			article.AddImage(i)
		}

		if article.Crawl.ID == 0 {
			article.Crawl = entity.Crawl{Next: time.Now()}
		}

		tx.Save(&article)
	}
	tx.Commit()

	return nil
}

func (r GormArticleRepository) Save(article *entity.Article) error {
	if result := r.db.Save(&article); result.Error != nil {
		return result.Error
	}

	return nil
}
