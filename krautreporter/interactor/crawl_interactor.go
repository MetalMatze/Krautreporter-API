package interactor

import (
	"github.com/metalmatze/krautreporter-api/krautreporter/entity"
)

type (
	crawlRepository interface {
		FindOutdatedAuthors() ([]*entity.Author, error)
		FindOutdatedArticles() ([]*entity.Article, error)
	}
	crawlAuthorRepository interface {
		FindByID(int) (*entity.Author, error)
		Save(*entity.Author) error
		SaveAll([]*entity.Author) error
	}

	crawlArticleRepository interface {
		FindByID(int) (*entity.Article, error)
		Save(*entity.Article) error
		SaveAll([]*entity.Article) error
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

// Authors

func (i *CrawlInteractor) FindOutdatedAuthors() ([]*entity.Author, error) {
	return i.crawlRepository.FindOutdatedAuthors()
}

func (i *CrawlInteractor) FindAuthorByID(id int) (*entity.Author, error) {
	return i.authorRepository.FindByID(id)
}

func (i *CrawlInteractor) SaveAuthor(a *entity.Author) error {
	return i.authorRepository.Save(a)
}

func (i *CrawlInteractor) SaveAuthors(as []*entity.Author) error {
	return i.authorRepository.SaveAll(as)
}

// Articles

func (i *CrawlInteractor) FindOutdatedArticles() ([]*entity.Article, error) {
	return i.crawlRepository.FindOutdatedArticles()
}

func (i *CrawlInteractor) FindArticleByID(id int) (*entity.Article, error) {
	return i.articleRepository.FindByID(id)
}

func (i *CrawlInteractor) SaveArticles(as []*entity.Article) error {
	return i.articleRepository.SaveAll(as)
}

func (i *CrawlInteractor) SaveArticle(a *entity.Article) error {
	return i.articleRepository.Save(a)
}
