package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/metalmatze/krautreporter-api/http/controller"
	"github.com/metalmatze/krautreporter-api/krautreporter"
)

func Routes(kr *krautreporter.Krautreporter, r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "hi")
	})

	c := controller.New(kr.Logger, kr.HTTPInteractor)

	authorsController := controller.AuthorsController{Controller: c}
	r.GET("/authors", authorsController.GetAuthors)
	r.GET("/authors/:id", authorsController.GetAuthor)

	articlesController := controller.ArticlesController{Controller: c}
	r.GET("/articles", articlesController.GetArticles)
	r.GET("/articles/:id", articlesController.GetArticle)

	crawlsController := controller.CrawlsController{Controller: c}
	r.GET("/crawls", crawlsController.GetCrawls)
}
