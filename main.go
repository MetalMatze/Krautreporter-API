package main

import (
	"log"

	"github.com/MetalMatze/Krautreporter-API/commands"
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
		AuthorRepository: repository.NewGormAuthorsRepository(gollection.DB),
	}

	gollection.AddCommands(commands.CrawlCommand(authorInteractor))
	gollection.AddRoutes(http.Routes(authorInteractor))

	if err := gollection.Run(); err != nil {
		log.Fatal(err)
	}
}
