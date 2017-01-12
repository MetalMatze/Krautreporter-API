package repository

import (
	"errors"
	"fmt"
	"strings"
	"time"

	krautreporter "github.com/metalmatze/krautreporter-api"
)

// MaxArticleID is used for sorting descending
const MaxArticleID int = 1234567890

// ErrArticleNotFound is returned if an article is not found by id
var ErrArticleNotFound = errors.New("Article not found")

// FindArticlesOlderThan returns a slice of Article that are older than an ID
func (r Repository) FindArticlesOlderThan(id int, number int) ([]*krautreporter.Article, error) {
	ordering := MaxArticleID
	if id != MaxArticleID {
		a, err := r.FindArticleByID(id)
		if err != nil {
			return nil, err
		}

		ordering = a.Ordering
	}

	var articles []*krautreporter.Article
	if result := r.DB.
		Preload("Images").
		Where("ordering < ?", ordering).
		Limit(number).
		Order("ordering desc").
		Find(&articles); result.Error != nil {
		return nil, result.Error
	}

	return articles, nil
}

// FindArticleByID returns an Article for the ID matching the parameter
func (r Repository) FindArticleByID(id int) (*krautreporter.Article, error) {
	if cached, exists := r.Cache.Get(fmt.Sprintf("articles.%d", id)); exists {
		return cached.(*krautreporter.Article), nil
	}

	var a krautreporter.Article
	r.DB.Preload("Images").Preload("Crawl").First(&a, "id = ?", id)

	if a.ID == 0 {
		return nil, ErrArticleNotFound
	}

	r.Cache.Set(fmt.Sprintf("authors.%d", a.ID), &a, 5*time.Second)

	return &a, nil
}

// SaveAllArticles takes a slice of Article and saves them to the database
func (r Repository) SaveAllArticles(articles []*krautreporter.Article) error {
	tx := r.DB.Begin()
	for i, a := range articles {
		article := krautreporter.Article{ID: a.ID}
		tx.Preload("Images").Preload("Crawl").FirstOrCreate(&article)

		article.Ordering = len(articles) - 1 - i
		article.Title = a.Title
		article.URL = a.URL
		article.Preview = a.Preview

		author := krautreporter.Author{}
		tx.First(&author, "name = ?", strings.TrimSpace(a.Author.Name))
		if author.ID == 0 {
			r.Logger.Log("msg", "Can't find author for article ", "author", a.Author.Name, "article", a.URL)
			continue
		}
		article.AuthorID = author.ID

		for _, i := range a.Images {
			article.AddImage(i)
		}

		if article.Crawl.ID == 0 {
			article.Crawl = &krautreporter.Crawl{Next: time.Now()}
		}

		tx.Save(&article)
	}
	tx.Commit()

	return nil
}

// SaveArticle takes an Article and saves it to the database
func (r Repository) SaveArticle(article *krautreporter.Article) error {
	if result := r.DB.Save(&article); result.Error != nil {
		return result.Error
	}

	return nil
}
