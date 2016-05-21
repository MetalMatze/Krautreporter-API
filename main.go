package main

import (
	"github.com/MetalMatze/Krautreporter-API/cli"
	"github.com/MetalMatze/Krautreporter-API/domain"
	"github.com/MetalMatze/Krautreporter-API/http"
	"github.com/MetalMatze/gollection"
	"github.com/MetalMatze/gollection/cache"
	"github.com/MetalMatze/gollection/database"
	"github.com/MetalMatze/gollection/router"
)

func main() {
	config := GetConfig()
	g := gollection.New(config)

	g.AddDB(database.Postgres(config))
	g.AddCache(cache.NewInMemory())
	g.AddRouter(router.NewGin())

	kr := domain.NewKrautreporter(g)

	g.AddRoutes(http.Routes(g, kr))
	g.AddCommands(cli.CrawlCommand(kr))

	if err := g.Run(); err != nil {
		g.Log.Crit("gollection error", "err", err)
	}
}
