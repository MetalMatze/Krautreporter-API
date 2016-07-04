package main

import (
	"os"

	"github.com/MetalMatze/Krautreporter-API/cli"
	"github.com/MetalMatze/Krautreporter-API/cmd"
	"github.com/MetalMatze/Krautreporter-API/config"
	"github.com/MetalMatze/Krautreporter-API/krautreporter"
	"github.com/go-kit/kit/log"
	"github.com/gollection/gollection"
)

func main() {
	logger := log.NewLogfmtLogger(os.Stderr)
	logger = log.NewContext(logger).With("ts", log.DefaultTimestampUTC)

	c := config.Config()
	g := gollection.New(logger, c.Config)

	gorm := cmd.Gorm(g, c)
	cache := cmd.Cache()

	kr := krautreporter.New(logger, gorm, cache)

	g.Cli.Commands = append(g.Cli.Commands, cli.SyncCommand(kr))
	g.Cli.Commands = append(g.Cli.Commands, cli.CrawlCommand(kr, logger))

	if err := g.Run(); err != nil {
		g.Logger.Log("msg", "Error running gollection", "err", err)
	}
}
