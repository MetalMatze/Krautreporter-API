package controllers

import (
	"net/http"
	"strconv"

	"github.com/MetalMatze/Krautreporter-API/domain/interactor"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	AuthorInteractor interactor.AuthorInteractor
}

func (controller *Controller) GetAuthors(c *gin.Context) {
	authors, err := controller.AuthorInteractor.GetAll()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(200, authors)
}

func (controller *Controller) GetAuthor(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	author, err := controller.AuthorInteractor.FindByID(id)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(200, author)
}
