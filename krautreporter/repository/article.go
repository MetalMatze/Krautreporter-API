package repository

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/MetalMatze/Krautreporter-API/krautreporter/entity"
	"github.com/jinzhu/gorm"
)

const MaxArticleID int = 1234567890

var ErrArticleNotFound = errors.New("Article not found")

type GormArticleRepository struct {
	DB *gorm.DB
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
	if result := r.DB.Where("ordering < ?", ordering).Limit(number).Order("ordering desc").Find(&articles); result.Error != nil {
		return nil, result.Error
	}

	return articles, nil
}

func (r GormArticleRepository) FindByID(id int) (*entity.Article, error) {
	var a entity.Article
	r.DB.First(&a, "id = ?", id)

	if a.ID == 0 {
		return nil, ErrArticleNotFound
	}

	return &a, nil
}

func (r GormArticleRepository) SaveAll(articles []entity.Article) error {
	tx := r.DB.Begin()
	for i, a := range articles {
		a.Ordering = len(articles) - 1 - i

		var author entity.Author
		r.DB.First(&author, "name = ?", strings.TrimSpace(a.Author.Name))
		if author.ID == 0 {
			log.Printf(`Can't find author "%s" for article %s`, a.Author.Name, a.URL)
			continue
		}
		a.AuthorID = author.ID
		a.Author = nil

		if _, err := r.FindByID(a.ID); err == ErrArticleNotFound {
			a.Crawl = entity.Crawl{Next: time.Now()} // Create crawl if article is new
			if result := tx.Create(&a); result.Error != nil {
				return result.Error
			}
		} else {
			if result := tx.Save(&a); result.Error != nil {
				return result.Error
			}
		}
	}
	tx.Commit()

	return nil
}

func (r GormArticleRepository) Save(article entity.Article) error {
	if result := r.DB.Save(&article); result.Error != nil {
		return result.Error
	}

	return nil
}
