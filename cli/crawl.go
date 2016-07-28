package cli

import (
	"strconv"
	"time"

	"github.com/MetalMatze/Krautreporter-API/krautreporter"
	"github.com/MetalMatze/Krautreporter-API/krautreporter/entity"
	"github.com/MetalMatze/Krautreporter-API/krautreporter/service"
	"github.com/fortytw2/radish"
	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/urfave/cli"
)

type CrawlInteractor interface {
	FindOutdatedAuthors() ([]*entity.Author, error)
	FindAuthorByID(int) (*entity.Author, error)
	SaveAuthor(*entity.Author) error
	SaveAuthors([]*entity.Author) error

	FindOutdatedArticles() ([]*entity.Article, error)
	FindArticleByID(int) (*entity.Article, error)
	SaveArticle(*entity.Article) error
	SaveArticles([]*entity.Article) error
}

const articleQueue = "articleQueue"
const authorQueue = "authorQueue"

func CrawlCommand(kr *krautreporter.Krautreporter, logger log.Logger, syncCounter prometheus.Counter, crawlCounter *prometheus.CounterVec) cli.Command {
	return cli.Command{
		Name:  "crawl",
		Usage: "Crawl article & authors",
		Action: func(c *cli.Context) (err error) {
			crawler := &crawler{
				interactor:   kr.CrawlInteractor,
				logger:       logger,
				crawlCounter: crawlCounter,
			}

			logger.Log("msg", "creating broker")
			broker := radish.NewMemBroker()

			logger.Log("msg", "creating pools")
			authorsPool := radish.NewPool(broker, authorQueue, crawler.authors, logger)
			articlesPool := radish.NewPool(broker, articleQueue, crawler.articles, logger)

			logger.Log("msg", "creating publishers")
			authorsPub, err := broker.Publisher(authorQueue)
			if err != nil {
				return err
			}

			articlesPub, err := broker.Publisher(articleQueue)
			if err != nil {
				return err
			}

			logger.Log("msg", "Adding workers")
			err = authorsPool.AddWorkers(10)
			if err != nil {
				return err
			}
			err = articlesPool.AddWorkers(10)
			if err != nil {
				return err
			}

			for {
				syncAuthor(kr.Logger, kr.CrawlInteractor)
				syncArticle(kr.Logger, kr.CrawlInteractor)
				syncCounter.Inc()

				authors, err := kr.CrawlInteractor.FindOutdatedAuthors()
				if err != nil {
					logger.Log("msg", "Can't get outdated authors from CrawlInteractor", "err", err)
				}

				for _, a := range authors {
					authorsPub.Publish([]byte(strconv.Itoa(a.ID)))
				}

				articles, err := kr.CrawlInteractor.FindOutdatedArticles()
				if err != nil {
					logger.Log("msg", "Can't get outdated articles from CrawlInteractor", "err", err)
				}

				for _, a := range articles {
					articlesPub.Publish([]byte(strconv.Itoa(a.ID)))
				}

				time.Sleep(5 * time.Minute)
			}

			logger.Log("msg", "Stopping the pool")
			err = authorsPool.Stop()
			if err != nil {
				return err
			}

			return nil
		},
	}
}

type crawler struct {
	interactor   CrawlInteractor
	logger       log.Logger
	crawlCounter *prometheus.CounterVec
}

func (c *crawler) authors(in []byte) ([][]byte, error) {
	start := time.Now()

	id, err := strconv.Atoi(string(in))
	if err != nil {
		return nil, err
	}

	a, err := c.interactor.FindAuthorByID(id)
	if err != nil {
		return nil, err
	}

	err = service.CrawlAuthor(a)
	if err != nil {
		c.crawlCounter.WithLabelValues("authors", "error").Inc()
		c.logger.Log("msg", "Crawling author has failed", "id", a.ID, "url", a.URL)
		return nil, err
	}

	err = c.interactor.SaveAuthor(a)
	if err != nil {
		c.crawlCounter.WithLabelValues("authors", "error").Inc()
		c.logger.Log("msg", "Failed to save crawled author", "id", a.ID, "url", a.URL, "duration", time.Since(start))
		return nil, err
	}

	c.crawlCounter.WithLabelValues("authors", "success").Inc()
	c.logger.Log("msg", "Author crawled successfully", "id", a.ID, "url", a.URL, "duration", time.Since(start))

	return nil, nil
}

func (c *crawler) articles(in []byte) ([][]byte, error) {
	start := time.Now()

	id, err := strconv.Atoi(string(in))
	if err != nil {
		return nil, err
	}

	a, err := c.interactor.FindArticleByID(id)
	if err != nil {
		return nil, err
	}

	err = service.CrawlArticle(a)
	if err != nil {
		c.crawlCounter.WithLabelValues("articles", "error").Inc()
		c.logger.Log("msg", "Crawling articles has failed", "id", a.ID, "url", a.URL)
		return nil, err
	}

	err = c.interactor.SaveArticle(a)
	if err != nil {
		c.crawlCounter.WithLabelValues("articles", "error").Inc()
		c.logger.Log("msg", "Failed to save crawled article", "id", a.ID, "url", a.URL, "duration", time.Since(start))
		return nil, err
	}

	c.crawlCounter.WithLabelValues("articles", "success").Inc()
	c.logger.Log("msg", "Article crawled successfully", "id", a.ID, "url", a.URL, "duration", time.Since(start))

	return nil, nil
}
