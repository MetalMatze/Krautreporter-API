package interactor

import (
	"github.com/MetalMatze/Krautreporter-API/krautreporter/entity"
)

type (
	crawlRepository interface {
		FindOutdatedAuthors() ([]entity.Author, error)
		FindOutdatedArticles() ([]entity.Article, error)
	}
	crawlAuthorRepository interface {
		Save(entity.Author) error
		SaveAll([]entity.Author) error
	}

	crawlArticleRepository interface {
		Save(entity.Article) error
		SaveAll([]entity.Article) error
	}
)

type CrawlInteractor struct {
	authorRepository  crawlAuthorRepository
	articleRepository crawlArticleRepository
	crawlRepository   crawlRepository
}

func NewCrawlInteractor(aur crawlAuthorRepository, arr crawlArticleRepository, cr crawlRepository) *CrawlInteractor {
	return &CrawlInteractor{
		authorRepository:  aur,
		articleRepository: arr,
		crawlRepository:   cr,
	}
}

func (i *CrawlInteractor) FindOutdatedAuthors() ([]entity.Author, error) {
	return i.crawlRepository.FindOutdatedAuthors()
}

func (i *CrawlInteractor) FindOutdatedArticles() ([]entity.Article, error) {
	return i.crawlRepository.FindOutdatedArticles()
}

func (i *CrawlInteractor) SaveAuthors(as []entity.Author) error {
	return i.authorRepository.SaveAll(as)
}

func (i *CrawlInteractor) SaveArticles(as []entity.Article) error {
	return i.articleRepository.SaveAll(as)
}

func (i *CrawlInteractor) SaveAuthor(a entity.Author) error {
	return i.authorRepository.Save(a)
}

func (i *CrawlInteractor) SaveArticle(a entity.Article) error {
	return i.articleRepository.Save(a)
}
