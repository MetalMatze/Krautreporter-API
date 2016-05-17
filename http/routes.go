package http

import (
	"github.com/MetalMatze/Krautreporter-API/domain"
	"github.com/MetalMatze/Krautreporter-API/http/controller"
	"github.com/MetalMatze/gollection/router"
)

func Routes(kr *domain.Krautreporter) func(router.Router) {
	return func(r router.Router) {
		authorsController := controller.AuthorsController{AuthorInteractor: kr.AuthorInteractor}
		r.GET("/authors", authorsController.GetAuthors)
		r.GET("/authors/:id", authorsController.GetAuthor)

		articlesController := controller.ArticlesController{ArticleInteractor: kr.ArticleInteractor}
		r.GET("/articles/:id", articlesController.GetArticle)
	}
}
