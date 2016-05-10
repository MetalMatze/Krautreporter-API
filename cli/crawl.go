package cli

import (
	"log"

	"github.com/MetalMatze/Krautreporter-API/domain/entity"
	"github.com/MetalMatze/Krautreporter-API/domain/service"
	"github.com/codegangsta/cli"
)

type AuthorInteractor interface {
	SaveAll(authors []entity.Author) error
}

type ArticleInteractor interface {
	SaveAll(authors []entity.Article) error
}

func CrawlCommand(authorInteractor AuthorInteractor, articleInteractor ArticleInteractor) cli.Command {
	return cli.Command{
		Name:  "crawl",
		Usage: "Display an inspiring quote",
		Action: func(c *cli.Context) {
			crawlAuthor(authorInteractor)
			crawlArticle(articleInteractor)
		},
		Subcommands: []cli.Command{{
			Name:  "authors",
			Usage: "Crawl all authors from krautreporter.de",
			Action: func(c *cli.Context) {
				crawlAuthor(authorInteractor)
			},
		}, {
			Name:  "articles",
			Usage: "Crawl all articles from krautreporter.de",
			Action: func(c *cli.Context) {
				crawlArticle(articleInteractor)
			},
		}},
	}
}

func crawlAuthor(authorInteractor AuthorInteractor) {
	authors, err := service.CrawlAuthor()
	if err != nil {
		log.Fatal(err)
	}

	if err := authorInteractor.SaveAll(authors); err != nil {
		log.Fatal(err)
	}
}

func crawlArticle(articlesInteractor ArticleInteractor) {
	articles, err := service.CrawlArticles()
	if err != nil {
		log.Fatal(err)
	}

	if err := articlesInteractor.SaveAll(articles); err != nil {
		log.Fatal(err)
	}
}
