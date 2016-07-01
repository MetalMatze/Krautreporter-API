package controller

import (
	"net/http"
	"strconv"

	"github.com/MetalMatze/Krautreporter-API/http/marshaller"
	"github.com/MetalMatze/Krautreporter-API/krautreporter/repository"
	"github.com/gin-gonic/gin"
)

type AuthorsController struct {
	*Controller
}

func (c *AuthorsController) GetAuthors(ctx *gin.Context) {
	authors, err := c.interactor.AllAuthors()
	if err != nil {
		c.logger.Log("Can't get all authors", "err", err.Error())
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}

	ctx.JSON(http.StatusOK, marshaller.Authors(authors))
}

func (c *AuthorsController) GetAuthor(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		c.logger.Log("Can't convert author id to int", "err", err.Error())
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}

	author, err := c.interactor.AuthorByID(id)
	//author, err := c.interactor.AuthorByID(id)
	if err != nil {
		if err == repository.ErrAuthorNotFound {
			c.logger.Log("Can't find author", "id", id)
			status := http.StatusNotFound
			ctx.JSON(status, gin.H{"message": http.StatusText(status), "status_code": status})
		}

		c.logger.Log("Can't get author", "id", id, "err", err.Error())
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}

	ctx.JSON(http.StatusOK, marshaller.Author(author))
}
