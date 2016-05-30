package interactor

import (
	"github.com/MetalMatze/Krautreporter-API/domain/entity"
)

type AuthorInteractor struct {
	AuthorRepository entity.AuthorRepository
}

func (i AuthorInteractor) GetAll() ([]*entity.Author, error) {
	return i.AuthorRepository.Find()
}

func (i AuthorInteractor) SaveAll(authors []entity.Author) error {
	return i.AuthorRepository.SaveAll(authors)
}

func (i AuthorInteractor) FindByID(id int) (*entity.Author, error) {
	return i.AuthorRepository.FindByID(id)
}

func (i AuthorInteractor) Save(author entity.Author) error {
	return i.AuthorRepository.Save(author)
}
