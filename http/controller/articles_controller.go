package controller

import (
	"net/http"
	"strconv"

	"github.com/MetalMatze/Krautreporter-API/domain/entity"
	"github.com/MetalMatze/gollection/router"
)

type ArticleInteractor interface {
	FindOlderThan(id int, number int) ([]*entity.Article, error)
	FindByID(id int) (*entity.Article, error)
}

type ArticlesController struct {
	ArticleInteractor ArticleInteractor
}

func (c *ArticlesController) GetArticles(req router.Request, res router.Response) error {
	id := 123456789
	if req.Query("olderthan") != "" {
		olderthan, err := strconv.Atoi(req.Query("olderthan"))
		if err != nil {
			return res.AbortWithStatus(http.StatusInternalServerError)
		}
		id = olderthan
	}

	articles, err := c.ArticleInteractor.FindOlderThan(id, 20)
	if err != nil {
		return res.AbortWithStatus(http.StatusInternalServerError)
	}

	return res.JSON(http.StatusOK, articles)
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
