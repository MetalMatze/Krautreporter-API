package http

import (
	"net/http"

	"github.com/MetalMatze/Krautreporter-API/http/controller"
	"github.com/MetalMatze/Krautreporter-API/krautreporter"
	"github.com/gollection/gollection"
	"github.com/gollection/gollection/router"
)

func Routes(g *gollection.Gollection, kr *krautreporter.Krautreporter) func(router.Router) {
	return func(r router.Router) {
		r.GET("/", func(res router.Response, req router.Request) error {
			return res.String(http.StatusOK, "hi")
		})

		r.GET("/health", func(res router.Response, req router.Request) error {
			if g.DB.DB().Ping() != nil {
				status := http.StatusInternalServerError
				return res.String(status, http.StatusText(status))
			}

			status := http.StatusOK
			return res.String(status, http.StatusText(status))
		})

		c := controller.Controller{}

		authorsController := controller.AuthorsController{Controller: c, AuthorInteractor: kr.AuthorInteractor}
		r.GET("/authors", authorsController.GetAuthors)
		r.GET("/authors/:id", authorsController.GetAuthor)

		articlesController := controller.ArticlesController{Controller: c, ArticleInteractor: kr.ArticleInteractor}
		r.GET("/articles", articlesController.GetArticles)
		r.GET("/articles/:id", articlesController.GetArticle)
	}
}
