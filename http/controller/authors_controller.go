package controller

import (
	"net/http"
	"strconv"

	"github.com/MetalMatze/Krautreporter-API/domain/entity"
	"github.com/MetalMatze/Krautreporter-API/domain/repository"
	"github.com/gollection/gollection/log"
	"github.com/gollection/gollection/router"
	"github.com/gin-gonic/gin"
)

type AuthorInteractor interface {
	GetAll() ([]*entity.Author, error)
	FindByID(id int) (*entity.Author, error)
}

type AuthorsController struct {
	AuthorInteractor AuthorInteractor
	Log              log.Logger
}

func (c *AuthorsController) GetAuthors(res router.Response, req router.Request) error {
	authors, err := c.AuthorInteractor.GetAll()
	if err != nil {
		c.Log.Info("Can't get all authors", "err", err.Error())
		return res.AbortWithStatus(http.StatusInternalServerError)
	}

	return res.JSON(http.StatusOK, authors)
}

func (c *AuthorsController) GetAuthor(res router.Response, req router.Request) error {
	id, err := strconv.Atoi(req.Param("id"))
	if err != nil {
		c.Log.Info("Can't convert author id to int", "err", err.Error())
		return res.AbortWithStatus(http.StatusInternalServerError)
	}

	author, err := c.AuthorInteractor.FindByID(id)
	if err != nil {
		if err == repository.ErrAuthorNotFound {
			c.Log.Debug("Can't find author", "id", id)
			status := http.StatusNotFound
			return res.JSON(status, gin.H{"message": http.StatusText(status), "status_code": status})
		}

		c.Log.Warn("Can't get author", "id", id, "err", err.Error())
		return res.AbortWithStatus(http.StatusInternalServerError)
	}

	return res.JSON(http.StatusOK, author)
}
