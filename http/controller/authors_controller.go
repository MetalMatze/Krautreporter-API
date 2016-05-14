package controller

import (
	"net/http"
	"strconv"

	"github.com/MetalMatze/Krautreporter-API/domain/entity"
	"github.com/MetalMatze/Krautreporter-API/domain/repository"
	"github.com/MetalMatze/gollection/router"
	"github.com/gin-gonic/gin"
)

type AuthorInteractor interface {
	GetAll() ([]*entity.Author, error)
	FindByID(id int) (*entity.Author, error)
}

type AuthorsController struct {
	AuthorInteractor AuthorInteractor
}

func (c *AuthorsController) GetAuthors(req router.Request, res router.Response) {
	authors, err := c.AuthorInteractor.GetAll()
	if err != nil {
		res.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	res.JSON(http.StatusOK, authors)
}

func (c *AuthorsController) GetAuthor(req router.Request, res router.Response) {
	id, err := strconv.Atoi(req.Param("id"))
	if err != nil {
		res.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	author, err := c.AuthorInteractor.FindByID(id)
	if err != nil {
		if err == repository.ErrAuthorNotFound {
			status := http.StatusNotFound
			res.JSON(status, gin.H{"message": http.StatusText(status), "status_code": status})
			return
		}

		res.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	res.JSON(http.StatusOK, author)
}
