package cli

import (
	"time"

	"github.com/MetalMatze/Krautreporter-API/krautreporter"
	"github.com/MetalMatze/Krautreporter-API/krautreporter/entity"
	"github.com/MetalMatze/Krautreporter-API/krautreporter/service"
	"github.com/MetalMatze/Krautreporter-API/workerqueue"
	"github.com/gollection/gollection/log"
	"github.com/urfave/cli"
)

type CrawlInteractor interface {
	FindOutdatedAuthors() ([]entity.Author, error)
	FindOutdatedArticles() ([]entity.Article, error)
	SaveAuthor(entity.Author) error
	SaveArticle(entity.Article) error
}

func CrawlCommand(kr *krautreporter.Krautreporter) cli.Command {
	return cli.Command{
		Name:  "crawl",
		Usage: "Crawl article & authors",
		Action: func(c *cli.Context) error {
			for {
				syncAuthor(kr.Log, kr.AuthorInteractor)
				syncArticle(kr.Log, kr.ArticleInteractor)

				crawler := newCrawler(kr.Log, kr.CrawlInteractor)
				crawler.authors()
				crawler.articles()

				time.Sleep(5 * time.Minute)
			}

			return nil
		},
	}
}

type crawler struct {
	interactor CrawlInteractor
	log        log.Logger
	wq         workerqueue.WorkerQueue
}

func newCrawler(log log.Logger, ci CrawlInteractor) crawler {
	return crawler{
		interactor: ci,
		log:        log,
		wq:         workerqueue.New(10),
	}
}

func (c crawler) authors() {
	authors, err := c.interactor.FindOutdatedAuthors()
	if err != nil {
		c.log.Warn("Can't get outdated authors from CrawlInteractor", "err", err)
	}

	for _, a := range authors {
		c.wq.Push(func() {
			start := time.Now()
			a, err := service.CrawlAuthor(a)
			if err != nil {
				c.log.Error("Crawling author has failed", "id", a.ID, "url", a.URL)
				return
			}

			err = c.interactor.SaveAuthor(a)
			if err != nil {
				c.log.Warn("Failed to save crawled author", "id", a.ID, "url", a.URL, "duration", time.Since(start))
				return
			}

			c.log.Info("Author crawled successfully", "id", a.ID, "url", a.URL, "duration", time.Since(start))
		})
	}
}

func (c crawler) articles() {
	articles, err := c.interactor.FindOutdatedArticles()
	if err != nil {
		c.log.Warn("Can't get outdated articles from CrawlInteractor", "err", err)
	}

	for _, a := range articles {
		c.wq.Push(func() {
			start := time.Now()
			a, err := service.CrawlArticle(a)
			if err != nil {
				c.log.Error("Crawling articles has failed", "id", a.ID, "url", a.URL)
				return
			}

			err = c.interactor.SaveArticle(a)
			if err != nil {
				c.log.Warn("Failed to save crawled article", "id", a.ID, "url", a.URL, "duration", time.Since(start))
				return
			}

			c.log.Info("Article crawled successfully", "id", a.ID, "url", a.URL, "duration", time.Since(start))
		})
	}
}
