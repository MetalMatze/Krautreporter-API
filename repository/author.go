package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/metalmatze/krautreporter-api/entity"
)

// ErrAuthorNotFound is returned if an article is not found by id
var ErrAuthorNotFound = errors.New("Author not found")

// FindAuthors returns a slice of all Authors
func (r Repository) FindAuthors() ([]*entity.Author, error) {
	if cached, exists := r.Cache.Get("authors.list"); exists {
		return cached.([]*entity.Author), nil
	}

	var authors []*entity.Author

	r.DB.Preload("Images").Order("ordering desc").Find(&authors)

	r.Cache.Set("authors.list", authors, time.Minute)

	return authors, nil
}

// FindAuthorByID returns an Author for the ID matching the parameter
func (r Repository) FindAuthorByID(id int) (*entity.Author, error) {
	if cached, exists := r.Cache.Get(fmt.Sprintf("authors.%d", id)); exists {
		return cached.(*entity.Author), nil
	}

	var author entity.Author
	r.DB.Preload("Images").Preload("Crawl").First(&author, "id = ?", id)

	if author.ID == 0 {
		return nil, ErrAuthorNotFound
	}

	r.Cache.Set(fmt.Sprintf("authors.%d", author.ID), &author, time.Minute)

	return &author, nil
}

// SaveAllAuthors takes a slice of Author and saves them to the database
func (r Repository) SaveAllAuthors(authors []*entity.Author) error {
	tx := r.DB.Begin()
	for _, a := range authors {
		author := entity.Author{ID: a.ID}
		tx.Preload("Crawl").Preload("Images").FirstOrCreate(&author)

		author.Ordering = a.Ordering
		author.Name = a.Name
		author.Title = a.Title
		author.URL = a.URL

		for _, i := range a.Images {
			author.AddImage(i)
		}

		if author.Crawl.ID == 0 {
			author.Crawl = entity.Crawl{Next: time.Now()}
		}

		tx.Save(&author)
	}
	tx.Commit()

	return nil
}

// SaveAuthor takes an Author and saves it to the database
func (r Repository) SaveAuthor(author *entity.Author) error {
	if result := r.DB.Save(&author); result.Error != nil {
		return result.Error
	}

	return nil
}
