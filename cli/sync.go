package cli

import (
	"log"

	"github.com/MetalMatze/Krautreporter-API/domain"
	"github.com/MetalMatze/Krautreporter-API/domain/entity"
	"github.com/MetalMatze/Krautreporter-API/domain/service"
	"github.com/urfave/cli"
)

type AuthorInteractor interface {
	SaveAll(authors []entity.Author) error
}

type ArticleInteractor interface {
	SaveAll(authors []entity.Article) error
}

func SyncCommand(kr *domain.Krautreporter) cli.Command {
	return cli.Command{
		Name:  "sync",
		Usage: "Sync authors & article from krautreporter.de",
		Action: func(c *cli.Context) {
			syncAuthor(kr.AuthorInteractor)
			syncArticle(kr.ArticleInteractor)
		},
		Subcommands: []cli.Command{{
			Name:  "authors",
			Usage: "Sync all authors from krautreporter.de",
			Action: func(c *cli.Context) {
				syncAuthor(kr.AuthorInteractor)
			},
		}, {
			Name:  "articles",
			Usage: "Sync all articles from krautreporter.de",
			Action: func(c *cli.Context) {
				syncArticle(kr.ArticleInteractor)
			},
		}},
	}
}

func syncAuthor(authorInteractor AuthorInteractor) {
	authors, err := service.SyncAuthor()
	if err != nil {
		log.Fatal(err)
	}

	if err := authorInteractor.SaveAll(authors); err != nil {
		log.Fatal(err)
	}
}

func syncArticle(articlesInteractor ArticleInteractor) {
	articles, err := service.SyncArticles()
	if err != nil {
		log.Fatal(err)
	}

	if err := articlesInteractor.SaveAll(articles); err != nil {
		log.Fatal(err)
	}
}
