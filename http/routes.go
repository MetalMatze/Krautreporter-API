package http

import (
	"github.com/MetalMatze/Krautreporter-API/http/controller"
	"github.com/gin-gonic/gin"
)

func Routes(authorInteractor controller.AuthorInteractor, articleInteractor controller.ArticleInteractor) func(router *gin.Engine) {
	return func(router *gin.Engine) {

		authorsController := controller.AuthorsController{AuthorInteractor: authorInteractor}
		router.GET("/authors", authorsController.GetAuthors)
		router.GET("/authors/:id", authorsController.GetAuthor)

		articlesController := controller.ArticlesController{Interactor: articleInteractor}
		router.GET("/articles/:id", articlesController.GetArticle)
	}
}
