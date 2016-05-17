package cli

import (
	"log"

	"github.com/MetalMatze/Krautreporter-API/domain"
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

func CrawlCommand(kr *domain.Krautreporter) cli.Command {
	return cli.Command{
		Name:  "crawl",
		Usage: "Crawl krautreporter.de to get authors & articles",
		Action: func(c *cli.Context) {
			crawlAuthor(kr.AuthorInteractor)
			crawlArticle(kr.ArticleInteractor)
		},
		Subcommands: []cli.Command{{
			Name:  "authors",
			Usage: "Crawl all authors from krautreporter.de",
			Action: func(c *cli.Context) {
				crawlAuthor(kr.AuthorInteractor)
			},
		}, {
			Name:  "articles",
			Usage: "Crawl all articles from krautreporter.de",
			Action: func(c *cli.Context) {
				crawlArticle(kr.ArticleInteractor)
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
