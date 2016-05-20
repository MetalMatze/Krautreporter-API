package interactor

import "github.com/MetalMatze/Krautreporter-API/domain/entity"

type ArticleRepository interface {
	SaveAll([]entity.Article) error
	FindOlderThan(id int, number int) ([]*entity.Article, error)
	FindByID(int) (*entity.Article, error)
}

type ArticleInteractor struct {
	ArticleRepository ArticleRepository
}

func (i ArticleInteractor) SaveAll(authors []entity.Article) error {
	return i.ArticleRepository.SaveAll(authors)
}

func (i ArticleInteractor) FindOlderThan(id int, number int) ([]*entity.Article, error) {
	return i.ArticleRepository.FindOlderThan(id, number)
}

func (i ArticleInteractor) FindByID(id int) (*entity.Article, error) {
	return i.ArticleRepository.FindByID(id)
}
