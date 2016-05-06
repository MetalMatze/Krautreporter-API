package interactor

import (
	"github.com/MetalMatze/Krautreporter-API/domain/entity"
	"github.com/MetalMatze/Krautreporter-API/domain/repository"
)

type AuthorInteractor struct {
	AuthorRepository repository.GormAuthorRepository
}

func (i *AuthorInteractor) GetAll() ([]entity.Author, error) {
	return i.AuthorRepository.Find()
}

func (i *AuthorInteractor) SaveAll(authors []entity.Author) error {
	return i.AuthorRepository.SaveAll(authors)
}

func (i *AuthorInteractor) FindByID(id int) (entity.Author, error) {
	return i.AuthorRepository.FindByID(id)
}
