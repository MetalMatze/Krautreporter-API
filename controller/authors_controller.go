package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/metalmatze/krautreporter-api/marshaller"
	"github.com/metalmatze/krautreporter-api/repository"
)

func (ctrl *Controller) GetAuthors(c *gin.Context) {
	authors, err := ctrl.Repository.FindAuthors()
	if err != nil {
		ctrl.Logger.Log("Can't get all authors", "err", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	c.JSON(http.StatusOK, marshaller.Authors(authors))
}

func (ctrl *Controller) GetAuthor(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ctrl.Logger.Log("Can't convert author id to int", "err", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	author, err := ctrl.Repository.FindAuthorByID(id)
	if err != nil {
		if err == repository.ErrAuthorNotFound {
			ctrl.Logger.Log("Can't find author", "id", id)
			status := http.StatusNotFound
			c.JSON(status, gin.H{"message": http.StatusText(status), "status_code": status})
		}

		ctrl.Logger.Log("Can't get author", "id", id, "err", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	c.JSON(http.StatusOK, marshaller.Author(author))
}
