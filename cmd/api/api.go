package main

import (
	"github.com/MetalMatze/Krautreporter-API/cli"
	"github.com/MetalMatze/Krautreporter-API/config"
	"github.com/MetalMatze/Krautreporter-API/http"
	"github.com/MetalMatze/Krautreporter-API/krautreporter"
	"github.com/gollection/gollection"
	"github.com/gollection/gollection/cache"
	"github.com/gollection/gollection/database/postgres"
	"github.com/gollection/gollection/router"
)

func main() {
	config := config.GetConfig()
	g := gollection.New(config)

	g.AddDB(postgres.New(g.Config))
	g.AddCache(cache.NewInMemory())
	g.AddRouter(router.NewGin())

	kr := krautreporter.New(g)

	g.AddRoutes(http.Routes(g, kr))
	g.AddCommands(cli.SeedCommand(kr))

	if err := g.Run(); err != nil {
		g.Log.Crit("Error running gollection", "err", err)
	}
}
