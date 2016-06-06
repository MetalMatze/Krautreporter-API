package cli

import (
	"github.com/MetalMatze/Krautreporter-API/krautreporter"
	"github.com/MetalMatze/Krautreporter-API/krautreporter/entity"
	"github.com/MetalMatze/Krautreporter-API/krautreporter/service"
	"github.com/gollection/gollection/log"
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
			syncAuthor(kr.Log, kr.AuthorInteractor)
			syncArticle(kr.Log, kr.ArticleInteractor)
			return nil
		},
		Subcommands: []cli.Command{{
			Name:  "authors",
			Usage: "Sync all authors from krautreporter.de",
			Action: func(c *cli.Context) error {
				syncAuthor(kr.Log, kr.AuthorInteractor)
				return nil
			},
		}, {
			Name:  "articles",
			Usage: "Sync all articles from krautreporter.de",
			Action: func(c *cli.Context) error {
				syncArticle(kr.Log, kr.ArticleInteractor)
				return nil
			},
		}},
	}
}

func syncAuthor(log log.Logger, authorInteractor AuthorInteractor) {
	authors, err := service.SyncAuthor(log)
	if err != nil {
		log.Error("Failed to sync authors", "err", err)
	}

	if err := authorInteractor.SaveAll(authors); err != nil {
		log.Error("Failed to save authors", "err", err)
	}
}

func syncArticle(log log.Logger, articlesInteractor ArticleInteractor) {
	articles, err := service.SyncArticles(log)
	if err != nil {
		log.Error("Failed to sync articles", "err", err)
	}

	if err := articlesInteractor.SaveAll(articles); err != nil {
		log.Error("Failed to save articles", "err", err)
	}
}
