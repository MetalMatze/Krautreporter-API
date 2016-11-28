package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/log"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	ginprometheus "github.com/mcuadros/go-gin-prometheus"
	"github.com/metalmatze/krautreporter-api/controller"
	"github.com/metalmatze/krautreporter-api/repository"
	gocache "github.com/patrickmn/go-cache"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/urfave/cli"
)

type Config struct {
	Addr string
	DSN  string
}

func main() {
	logger := log.NewLogfmtLogger(os.Stderr)
	logger = log.NewContext(logger).With("ts", log.DefaultTimestampUTC)

	err := godotenv.Load()
	if err != nil {
		logger.Log("level", "fatal", "err", err)
	}

	config := &Config{
		Addr: os.Getenv("ADDR"),
		DSN:  os.Getenv("DSN"),
	}

	db, err := gorm.Open("postgres", config.DSN)
	if err != nil {
		panic(err)
	}

	app := cli.NewApp()

	app.Commands = []cli.Command{{
		Name:   "serve",
		Action: serve(logger, config, db),
	}}

	if err := app.Run(os.Args); err != nil {
		logger.Log("level", "fatal", "err", err)
	}

}

func serve(logger log.Logger, config *Config, db *gorm.DB) func(*cli.Context) error {
	return func(c *cli.Context) error {
		router := gin.New()
		router.Use(gin.Recovery())

		router.GET("/", func(c *gin.Context) {
			c.String(http.StatusOK, "hi")
		})

		router.GET("/health", func(c *gin.Context) {
			status := http.StatusOK
			c.String(status, http.StatusText(status))
		})

		repo := repository.Repository{
			Cache:  gocache.New(5*time.Minute, 30*time.Second),
			DB:     db,
			Logger: logger,
		}
		ctrl := controller.Controller{
			Logger:     logger,
			Repository: repo,
		}

		ginprom := ginprometheus.NewPrometheus("gin")
		ginprom.Use(router)

		router.GET("/authors", ctrl.GetAuthors)
		router.GET("/authors/:id", ctrl.GetAuthor)

		router.GET("/articles", ctrl.GetArticles)
		router.GET("/articles/:id", ctrl.GetArticle)

		router.GET("/crawls", ctrl.GetCrawls)

		go func() {
			http.Handle("/metrics", prometheus.Handler())
			if err := http.ListenAndServe(":8081", nil); err != nil {
				panic(err)
			}
		}()

		if err := router.Run(config.Addr); err != nil {
			return err
		}

		return nil
	}
}
