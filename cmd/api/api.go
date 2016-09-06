package main

import (
	nethttp "net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/log"
	"github.com/gollection/gollection"
	"github.com/metalmatze/krautreporter-api/cmd"
	"github.com/metalmatze/krautreporter-api/config"
	"github.com/metalmatze/krautreporter-api/http"
	"github.com/metalmatze/krautreporter-api/krautreporter"
)

func main() {
	logger := log.NewLogfmtLogger(os.Stderr)
	logger = log.NewContext(logger).With("ts", log.DefaultTimestampUTC)

	c := config.Config()
	g := gollection.New(logger, c.Config)

	router := cmd.Gin(g, c)
	gorm := cmd.Gorm(g, c)
	cache := cmd.Cache()

	kr := krautreporter.New(logger, gorm, cache)

	http.Routes(kr, router)

	router.GET("/health", func(c *gin.Context) {
		if gorm.DB().Ping() != nil {
			status := nethttp.StatusInternalServerError
			c.String(status, nethttp.StatusText(status))
		}
		status := nethttp.StatusOK
		c.String(status, nethttp.StatusText(status))
	})

	g.Cli.Commands = append(g.Cli.Commands)

	if err := g.Run(); err != nil {
		g.Logger.Log("msg", "Error running gollection")
	}
}
