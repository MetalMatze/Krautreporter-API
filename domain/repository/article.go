package repository

import (
	"errors"
	"log"
	"strings"

	"github.com/MetalMatze/Krautreporter-API/domain/entity"
	"github.com/jinzhu/gorm"
)

var ErrArticleNotFound = errors.New("Article not found")

type GormArticleRepository struct {
	DB *gorm.DB
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
