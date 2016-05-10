package controller

import (
	"net/http"
	"strconv"

	"github.com/MetalMatze/Krautreporter-API/domain/entity"
	"github.com/MetalMatze/Krautreporter-API/domain/repository"
	"github.com/gin-gonic/gin"
)

type AuthorInteractor interface {
	GetAll() ([]*entity.Author, error)
	FindByID(id int) (*entity.Author, error)
}

type AuthorsController struct {
	AuthorInteractor AuthorInteractor
}

func (controller *AuthorsController) GetAuthors(c *gin.Context) {
	authors, err := controller.AuthorInteractor.GetAll()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(200, authors)
}

func (controller *AuthorsController) GetAuthor(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	author, err := controller.AuthorInteractor.FindByID(id)
	if err != nil {
		if err == repository.ErrAuthorNotFound {
			status := http.StatusNotFound
			c.JSON(status, gin.H{"message": http.StatusText(status), "status_code": status})
			return
		}

		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(200, author)
}
