package controller

import (
	"net/http"
	"strconv"

	"github.com/MetalMatze/Krautreporter-API/domain/entity"
	"github.com/gin-gonic/gin"
)

type ArticleInteractor interface {
	FindByID(id int) (*entity.Article, error)
}

type ArticlesController struct {
	Interactor ArticleInteractor
}

func (controller *ArticlesController) GetArticle(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	article, err := controller.Interactor.FindByID(id)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, article)
}
