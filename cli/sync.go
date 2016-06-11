package cli

import (
	"github.com/MetalMatze/Krautreporter-API/krautreporter"
	"github.com/MetalMatze/Krautreporter-API/krautreporter/service"
	"github.com/gollection/gollection/log"
	"github.com/urfave/cli"
)

func SyncCommand(kr *krautreporter.Krautreporter) cli.Command {
	return cli.Command{
		Name:  "sync",
		Usage: "Sync authors & article from krautreporter.de",
		Action: func(c *cli.Context) error {
			syncAuthor(kr.Log, kr.CrawlInteractor)
			syncArticle(kr.Log, kr.CrawlInteractor)
			return nil
		},
		Subcommands: []cli.Command{{
			Name:  "authors",
			Usage: "Sync all authors from krautreporter.de",
			Action: func(c *cli.Context) error {
				syncAuthor(kr.Log, kr.CrawlInteractor)
				return nil
			},
		}, {
			Name:  "articles",
			Usage: "Sync all articles from krautreporter.de",
			Action: func(c *cli.Context) error {
				syncArticle(kr.Log, kr.CrawlInteractor)
				return nil
			},
		}},
	}
}

func syncAuthor(log log.Logger, ci CrawlInteractor) {
	authors, err := service.SyncAuthor(log)
	if err != nil {
		log.Error("Failed to sync authors", "err", err)
	}

	if err := ci.SaveAuthors(authors); err != nil {
		log.Error("Failed to save authors", "err", err)
	}
}

func syncArticle(log log.Logger, ci CrawlInteractor) {
	articles, err := service.SyncArticles(log)
	if err != nil {
		log.Error("Failed to sync articles", "err", err)
	}

	if err := ci.SaveArticles(articles); err != nil {
		log.Error("Failed to save articles", "err", err)
	}
}
