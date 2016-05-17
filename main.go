package main

import (
	"log"

	"github.com/MetalMatze/Krautreporter-API/cli"
	"github.com/MetalMatze/Krautreporter-API/domain"
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

	kr := domain.NewKrautreporter(g)

	g.AddRoutes(http.Routes(kr))
	g.AddCommands(cli.CrawlCommand(kr))

	if err := g.Run(); err != nil {
		log.Fatal(err)
	}
}
