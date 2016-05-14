package controller

import (
	"github.com/MetalMatze/Krautreporter-API/domain/entity"
	"github.com/MetalMatze/gollection/router"
	"net/http"
	"strconv"
)

type ArticleInteractor interface {
	FindByID(id int) (*entity.Article, error)
}

type ArticlesController struct {
	ArticleInteractor ArticleInteractor
}

func (c *ArticlesController) GetArticle(req router.Request, res router.Response) error {
	id, err := strconv.Atoi(req.Param("id"))
	if err != nil {
		return res.AbortWithStatus(http.StatusInternalServerError)
	}

	article, err := c.ArticleInteractor.FindByID(id)
	if err != nil {
		return res.AbortWithStatus(http.StatusNotFound)
	}

	return res.JSON(http.StatusOK, article)
}
