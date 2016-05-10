package main

import (
	"log"

	"github.com/MetalMatze/Krautreporter-API/cli"
	"github.com/MetalMatze/Krautreporter-API/domain/entity"
	"github.com/MetalMatze/Krautreporter-API/domain/interactor"
	"github.com/MetalMatze/Krautreporter-API/domain/repository"
	"github.com/MetalMatze/Krautreporter-API/http"
	"github.com/MetalMatze/gollection"
	"github.com/MetalMatze/gollection/database"
)

func main() {
	config := GetConfig() // Get all env variables and populate the config with it
	gollection := gollection.New(config)

	gollection.AddDB(database.Postgres(config))

	authorInteractor := interactor.AuthorInteractor{
		AuthorRepository: repository.GormAuthorRepository{DB: gollection.DB},
	}
	articlesInteractor := interactor.ArticleInteractor{
		ArticleRepository: repository.GormArticleRepository{DB: gollection.DB},
	}

	gollection.DB.AutoMigrate(entity.Author{}, entity.Article{})

	gollection.AddCommands(cli.CrawlCommand(authorInteractor, articlesInteractor))
	gollection.AddRoutes(http.Routes(authorInteractor, articlesInteractor))

	if err := gollection.Run(); err != nil {
		log.Fatal(err)
	}
}
