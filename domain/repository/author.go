package repository

import (
	"github.com/MetalMatze/Krautreporter-API/domain/entity"
	"github.com/jinzhu/gorm"
)

type GormRepository struct {
	DB *gorm.DB
}

type GormAuthorRepository GormRepository

func NewGormAuthorsRepository(db *gorm.DB) GormAuthorRepository {
	return GormAuthorRepository{
		DB: db,
	}
}

func (r GormAuthorRepository) Find() ([]entity.Author, error) {
	var authors []entity.Author
	r.DB.Order("ordering desc").Find(&authors)

	return authors, nil
}

func (r GormAuthorRepository) FindByID(id int) (entity.Author, error) {
	var author entity.Author
	r.DB.First(&author, "id = ?", id)

	return author, nil
}

func (r GormAuthorRepository) SaveAll(authors []entity.Author) error {
	for _, author := range authors {
		r.DB.Create(&author)
	}
	return nil
}
