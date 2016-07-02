package cmd

import (
	"fmt"
	"time"

	"github.com/MetalMatze/Krautreporter-API/config"
	"github.com/gin-gonic/gin"
	"github.com/gollection/gollection"
	"github.com/gollection/gollection/database/gorm/postgres"
	gogin "github.com/gollection/gollection/router/gin"
	"github.com/jinzhu/gorm"
	gocache "github.com/patrickmn/go-cache"
)

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
