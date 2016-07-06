package main

import (
	"net/http"
	"os"

	"github.com/MetalMatze/Krautreporter-API/cli"
	"github.com/MetalMatze/Krautreporter-API/cmd"
	"github.com/MetalMatze/Krautreporter-API/config"
	"github.com/MetalMatze/Krautreporter-API/krautreporter"
	"github.com/go-kit/kit/log"
	"github.com/gollection/gollection"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	syncCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "sync",
			Help: "Number of syncs",
		},
	)
	crawlCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "crawls",
			Help: "Number of crawls",
		},
		[]string{"type"},
	)
)

func init() {
	prometheus.MustRegister(syncCounter)
	prometheus.MustRegister(crawlCounter)
}

func main() {
	logger := log.NewLogfmtLogger(os.Stderr)
	logger = log.NewContext(logger).With("ts", log.DefaultTimestampUTC)

	c := config.Config()
	g := gollection.New(logger, c.Config)

	gorm := cmd.Gorm(g, c)
	cache := cmd.Cache()

	kr := krautreporter.New(logger, gorm, cache)

	g.Cli.Commands = append(g.Cli.Commands, cli.SyncCommand(kr))
	g.Cli.Commands = append(g.Cli.Commands, cli.CrawlCommand(kr, logger, syncCounter, crawlCounter))

	go metrics()

	if err := g.Run(); err != nil {
		g.Logger.Log("msg", "Error running gollection", "err", err)
	}
}

func metrics() {
	http.Handle("/metrics", prometheus.Handler())
	http.ListenAndServe(":8080", nil)
}
