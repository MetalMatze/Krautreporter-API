package commands

import (
	"log"

	"github.com/MetalMatze/Krautreporter-API/domain/interactor"
	"github.com/MetalMatze/Krautreporter-API/domain/service"
	"github.com/codegangsta/cli"
)

func CrawlCommand(authorInteractor interactor.AuthorInteractor, articleInteractor interactor.ArticleInteractor) cli.Command {
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

func crawlAuthor(authorInteractor interactor.AuthorInteractor) {
	authors, err := service.CrawlAuthor()
	if err != nil {
		log.Fatal(err)
	}

	if err := authorInteractor.SaveAll(authors); err != nil {
		log.Fatal(err)
	}
}

func crawlArticle(articlesInteractor interactor.ArticleInteractor) {
	articles, err := service.CrawlArticles()
	if err != nil {
		log.Fatal(err)
	}

	if err := articlesInteractor.SaveAll(articles); err != nil {
		log.Fatal(err)
	}
}
