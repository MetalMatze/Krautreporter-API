package interactor

import (
	"github.com/MetalMatze/Krautreporter-API/krautreporter/entity"
)

type CrawlRepository interface {
	FindOutdatedAuthors() ([]entity.Author, error)
}

type CrawlInteractor struct {
	AuthorRepository AuthorRepository
	CrawlRepository  CrawlRepository
}

func (i CrawlInteractor) FindOutdatedAuthors() ([]entity.Author, error) {
	return i.CrawlRepository.FindOutdatedAuthors()
}

func (i CrawlInteractor) SaveAuthor(a entity.Author) error {
	return i.AuthorRepository.Save(a)
}
