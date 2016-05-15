package repository

import (
	"errors"

	"github.com/MetalMatze/Krautreporter-API/domain/entity"
	"github.com/jinzhu/gorm"
)

var ErrAuthorNotFound = errors.New("Author not found")

type GormRepository struct {
	DB *gorm.DB
}

type GormAuthorRepository GormRepository

func (r GormAuthorRepository) Find() ([]*entity.Author, error) {
	var authors []*entity.Author
	r.DB.Preload("Images").Order("ordering desc").Find(&authors)

	return authors, nil
}

func (r GormAuthorRepository) FindByID(id int) (*entity.Author, error) {
	var author entity.Author
	r.DB.Preload("Images").First(&author, "id = ?", id)

	if author.ID == 0 {
		return nil, ErrAuthorNotFound
	}

	return &author, nil
}

func (r GormAuthorRepository) SaveAll(authors []entity.Author) error {
	tx := r.DB.Begin()
	for _, a := range authors {
		if author, err := r.FindByID(a.ID); err == ErrAuthorNotFound {
			if result := tx.Create(&a); result.Error != nil {
				return result.Error
			}
		} else {
			if len(author.Images) > 0 {
				a.Images = author.Images
			}

			if result := tx.Save(&a); result.Error != nil {
				return result.Error
			}
		}
	}
	tx.Commit()

	return nil
}
