package main

import (
	"log"

	"github.com/MetalMatze/Krautreporter-API/cli"
	"github.com/MetalMatze/Krautreporter-API/domain/interactor"
	"github.com/MetalMatze/Krautreporter-API/domain/repository"
	"github.com/MetalMatze/Krautreporter-API/http"
	"github.com/MetalMatze/gollection"
	"github.com/MetalMatze/gollection/database"
	"github.com/MetalMatze/gollection/router"
)

func main() {
	config := GetConfig()
	g := gollection.New(config)

	g.AddDB(database.Postgres(config))
	g.AddRedis(gollection.NewRedis(config))

	g.AddRouter(router.NewGin())

	authorInteractor := interactor.AuthorInteractor{
		AuthorRepository: repository.GormAuthorRepository{DB: g.DB},
	}
	articlesInteractor := interactor.ArticleInteractor{
		ArticleRepository: repository.GormArticleRepository{DB: g.DB},
	}

	krRoutes := http.KrautreporterRoutes{
		AuthorInteractor:  authorInteractor,
		ArticleInteractor: articlesInteractor,
	}
	g.AddRoutes(krRoutes.Routes)

	g.AddCommands(cli.CrawlCommand(authorInteractor, articlesInteractor))

	if err := g.Run(); err != nil {
		log.Fatal(err)
	}
}
