package main

import (
	"github.com/MetalMatze/Krautreporter-API/cli"
	"github.com/MetalMatze/Krautreporter-API/config"
	"github.com/MetalMatze/Krautreporter-API/domain"
	"github.com/MetalMatze/gollection"
	"github.com/MetalMatze/gollection/cache"
	"github.com/MetalMatze/gollection/router"
	"github.com/MetalMatze/gollection/database/postgres"
)

func main() {
	config := config.GetConfig()
	g := gollection.New(config)

	g.AddDB(postgres.New(g.Config))
	g.AddCache(cache.NewInMemory())
	g.AddRouter(router.NewGin())

	kr := domain.NewKrautreporter(g)

	g.AddCommands(
		cli.SyncCommand(kr),
		cli.CrawlCommand(kr),
	)

	if err := g.Run(); err != nil {
		g.Log.Crit("Error running gollection", "err", err)
	}
}
