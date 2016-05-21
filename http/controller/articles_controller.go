package controller

import (
	"net/http"
	"strconv"

	"github.com/MetalMatze/Krautreporter-API/domain/entity"
	"github.com/MetalMatze/Krautreporter-API/domain/repository"
	"github.com/MetalMatze/gollection/router"
	"github.com/gin-gonic/gin"
)

type ArticleInteractor interface {
	FindOlderThan(id int, number int) ([]*entity.Article, error)
	FindByID(id int) (*entity.Article, error)
}

type ArticlesController struct {
	ArticleInteractor ArticleInteractor
}

func (c *ArticlesController) GetArticles(req router.Request, res router.Response) error {
	id := repository.MaxArticleID
	if req.Query("olderthan") != "" {
		olderthan, err := strconv.Atoi(req.Query("olderthan"))
		if err != nil {
			return res.AbortWithStatus(http.StatusInternalServerError)
		}
		id = olderthan
	}

	articles, err := c.ArticleInteractor.FindOlderThan(id, 20)
	if err != nil {
		if err == repository.ErrArticleNotFound {
			status := http.StatusNotFound
			return res.JSON(status, gin.H{
				"message":     http.StatusText(status),
				"status_code": status,
			})
		}

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
