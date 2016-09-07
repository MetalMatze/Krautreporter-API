package interactor

import "github.com/metalmatze/krautreporter-api/entity"

type httpAuthorRepository interface {
	Find() ([]*entity.Author, error)
	FindByID(int) (*entity.Author, error)
}

type httpArticleRepository interface {
	FindOlderThan(id int, number int) ([]*entity.Article, error)
	FindByID(int) (*entity.Article, error)
}

type httpCrawlRepository interface {
	NextCrawls(limit int) ([]*entity.Crawl, error)
}

type HTTPInteractor struct {
	authorRepository  httpAuthorRepository
	articleRepository httpArticleRepository
	crawlRepository   httpCrawlRepository
}

func NewHTTPInteractor(aur httpAuthorRepository, arr httpArticleRepository, cr httpCrawlRepository) *HTTPInteractor {
	return &HTTPInteractor{
		authorRepository:  aur,
		articleRepository: arr,
		crawlRepository:   cr,
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

func (i HTTPInteractor) NextCrawls(limit int) ([]*entity.Crawl, error) {
	return i.crawlRepository.NextCrawls(limit)
}
