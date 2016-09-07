package main

import (
	"fmt"
	nethttp "net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/log"
	"github.com/gollection/gollection"
	"github.com/gollection/gollection/database/gorm/postgres"
	gogin "github.com/gollection/gollection/router/gin"
	"github.com/jinzhu/gorm"
	"github.com/metalmatze/krautreporter-api/config"
	"github.com/metalmatze/krautreporter-api/http"
	"github.com/metalmatze/krautreporter-api/krautreporter"
	gocache "github.com/patrickmn/go-cache"
	"github.com/prometheus/client_golang/prometheus"
)

func main() {
	logger := log.NewLogfmtLogger(os.Stderr)
	logger = log.NewContext(logger).With("ts", log.DefaultTimestampUTC)

	c := config.Config()
	g := gollection.New(logger, c.Config)

	router := Gin(g, c)
	gorm := Gorm(g, c)
	cache := Cache()

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

	go func() {
		nethttp.Handle("/metrics", prometheus.Handler())
		nethttp.ListenAndServe(":8080", nil)
	}()

	if err := g.Run(); err != nil {
		g.Logger.Log("msg", "Error running gollection")
	}
}

func Gorm(g *gollection.Gollection, c config.AppConfig) *gorm.DB {
	gorm, err := postgres.New(g.Logger, c.DatabaseConfig)
	if err != nil {
		fmt.Printf("%+v\n", c.DatabaseConfig)
		panic(err)
	}

	return gorm
}

func Gin(g *gollection.Gollection, c config.AppConfig) *gin.Engine {
	ginWrapper := gogin.New(g.Logger, c.RouterConfig)
	g.Register(ginWrapper)

	return ginWrapper.Engine
}

func Cache() *gocache.Cache {
	return gocache.New(5*time.Minute, 30*time.Second)
}
