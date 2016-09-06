package cli

import (
	"github.com/go-kit/kit/log"
	"github.com/metalmatze/krautreporter-api/krautreporter"
	"github.com/metalmatze/krautreporter-api/krautreporter/service"
	"github.com/urfave/cli"
)

func SyncCommand(kr *krautreporter.Krautreporter) cli.Command {
	return cli.Command{
		Name:  "sync",
		Usage: "Sync authors & article from krautreporter.de",
		Action: func(c *cli.Context) error {
			syncAuthor(kr.Logger, kr.CrawlInteractor)
			syncArticle(kr.Logger, kr.CrawlInteractor)
			return nil
		},
		Subcommands: []cli.Command{{
			Name:  "authors",
			Usage: "Sync all authors from krautreporter.de",
			Action: func(c *cli.Context) error {
				syncAuthor(kr.Logger, kr.CrawlInteractor)
				return nil
			},
		}, {
			Name:  "articles",
			Usage: "Sync all articles from krautreporter.de",
			Action: func(c *cli.Context) error {
				syncArticle(kr.Logger, kr.CrawlInteractor)
				return nil
			},
		}},
	}
}

func syncAuthor(logger log.Logger, ci CrawlInteractor) {
	authors, err := service.SyncAuthor(logger)
	if err != nil {
		logger.Log("msg", "Failed to sync authors", "err", err)
	}

	if err := ci.SaveAuthors(authors); err != nil {
		logger.Log("msg", "Failed to save authors", "err", err)
	}
}

func syncArticle(logger log.Logger, ci CrawlInteractor) {
	articles, err := service.SyncArticles(logger)
	if err != nil {
		logger.Log("msg", "Failed to sync articles", "err", err)
	}

	if err := ci.SaveArticles(articles); err != nil {
		logger.Log("msg", "Failed to save articles", "err", err)
	}
}
