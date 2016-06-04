package cli

import (
	"log"

	"github.com/MetalMatze/Krautreporter-API/krautreporter"
	"github.com/MetalMatze/Krautreporter-API/krautreporter/entity"
	"github.com/MetalMatze/Krautreporter-API/krautreporter/service"
	"github.com/urfave/cli"
)

type AuthorInteractor interface {
	SaveAll(authors []entity.Author) error
}

type ArticleInteractor interface {
	SaveAll(authors []entity.Article) error
}

func SyncCommand(kr *krautreporter.Krautreporter) cli.Command {
	return cli.Command{
		Name:  "sync",
		Usage: "Sync authors & article from krautreporter.de",
		Action: func(c *cli.Context) error {
			syncAuthor(kr.AuthorInteractor)
			syncArticle(kr.ArticleInteractor)
			return nil
		},
		Subcommands: []cli.Command{{
			Name:  "authors",
			Usage: "Sync all authors from krautreporter.de",
			Action: func(c *cli.Context) error {
				syncAuthor(kr.AuthorInteractor)
				return nil
			},
		}, {
			Name:  "articles",
			Usage: "Sync all articles from krautreporter.de",
			Action: func(c *cli.Context) error {
				syncArticle(kr.ArticleInteractor)
				return nil
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
