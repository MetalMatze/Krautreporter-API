package interactor

import "github.com/MetalMatze/Krautreporter-API/krautreporter/entity"

// ArticleRepository is an interface that ArticleInteractor needs to get articles
type ArticleRepository interface {
	SaveAll([]entity.Article) error
	FindOlderThan(id int, number int) ([]*entity.Article, error)
	FindByID(int) (*entity.Article, error)
	Save(entity.Article) error
}

// ArticleInteractor is a boundary to interact with articles
type ArticleInteractor struct {
	ArticleRepository ArticleRepository
}

// SaveAll gets a slice of Articles and saves them
func (i ArticleInteractor) SaveAll(authors []entity.Article) error {
	return i.ArticleRepository.SaveAll(authors)
}

// FindOlderThan takes an id and returns the following number articles
func (i ArticleInteractor) FindOlderThan(id int, number int) ([]*entity.Article, error) {
	return i.ArticleRepository.FindOlderThan(id, number)
}

// FindByID takes an id an returns an Article
func (i ArticleInteractor) FindByID(id int) (*entity.Article, error) {
	return i.ArticleRepository.FindByID(id)
}
