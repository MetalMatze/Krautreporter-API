package http

import (
	"github.com/MetalMatze/Krautreporter-API/domain/interactor"
	"github.com/MetalMatze/Krautreporter-API/http/controller"
	"github.com/gin-gonic/gin"
)

func Routes(authorInteractor interactor.AuthorInteractor) func(router *gin.Engine) {
	return func(router *gin.Engine) {

		authorsController := controller.AuthorsController{AuthorInteractor: authorInteractor}
		router.GET("/authors", authorsController.GetAuthors)
		router.GET("/authors/:id", authorsController.GetAuthor)
	}
}
