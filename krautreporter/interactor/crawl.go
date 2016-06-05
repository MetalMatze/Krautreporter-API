package interactor

import (
	"github.com/MetalMatze/Krautreporter-API/krautreporter/entity"
)

type CrawlRepository interface {
	FindOutdatedAuthors() ([]entity.Author, error)
	FindOutdatedArticles() ([]entity.Article, error)
}

type CrawlInteractor struct {
	AuthorRepository  AuthorRepository
	ArticleRepository ArticleRepository
	CrawlRepository   CrawlRepository
}

func (i CrawlInteractor) FindOutdatedAuthors() ([]entity.Author, error) {
	return i.CrawlRepository.FindOutdatedAuthors()
}

func (i CrawlInteractor) FindOutdatedArticles() ([]entity.Article, error) {
	return i.CrawlRepository.FindOutdatedArticles()
}

func (i CrawlInteractor) SaveAuthor(a entity.Author) error {
	return i.AuthorRepository.Save(a)
}

func (i CrawlInteractor) SaveArticle(a entity.Article) error {
	return i.ArticleRepository.Save(a)
}
