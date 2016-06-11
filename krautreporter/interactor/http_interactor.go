package interactor

import "github.com/MetalMatze/Krautreporter-API/krautreporter/entity"

type httpAuthorRepository interface {
	Find() ([]*entity.Author, error)
	FindByID(int) (*entity.Author, error)
}

type httpArticleRepository interface {
	FindOlderThan(id int, number int) ([]*entity.Article, error)
	FindByID(int) (*entity.Article, error)
}

type HTTPInteractor struct {
	authorRepository  httpAuthorRepository
	articleRepository httpArticleRepository
}

func NewHTTPInteractor(aur httpAuthorRepository, arr httpArticleRepository) *HTTPInteractor {
	return &HTTPInteractor{
		authorRepository:  aur,
		articleRepository: arr,
	}
}

func (i HTTPInteractor) AllAuthors() ([]*entity.Author, error) {
	return i.authorRepository.Find()
}

func (i HTTPInteractor) AuthorByID(id int) (*entity.Author, error) {
	return i.authorRepository.FindByID(id)
}

func (i HTTPInteractor) ArticlesOlderThan(id int, number int) ([]*entity.Article, error) {
	return i.articleRepository.FindOlderThan(id, number)
}

func (i HTTPInteractor) ArticleByID(id int) (*entity.Article, error) {
	return i.articleRepository.FindByID(id)
}
