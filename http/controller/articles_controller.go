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

type ArticleInteractor interface {
	FindOlderThan(id int, number int) ([]*entity.Article, error)
	FindByID(id int) (*entity.Article, error)
}

type ArticlesController struct {
	ArticleInteractor ArticleInteractor
	Log               log.Logger
}

func (c *ArticlesController) GetArticles(res router.Response, req router.Request) error {
	id := repository.MaxArticleID
	if req.Query("olderthan") != "" {
		olderthan, err := strconv.Atoi(req.Query("olderthan"))
		if err != nil {
			c.Log.Info("Can't convert olderthan id to int", "err", err.Error())
			return res.AbortWithStatus(http.StatusInternalServerError)
		}
		id = olderthan
	}

	articles, err := c.ArticleInteractor.FindOlderThan(id, 20)
	if err != nil {
		if err == repository.ErrArticleNotFound {
			c.Log.Debug("Can't find olderthan article", "id", id)
			status := http.StatusNotFound
			return res.JSON(status, gin.H{
				"message":     http.StatusText(status),
				"status_code": status,
			})
		}

		c.Log.Warn("Failed to get olderthan articles", "id", id, "err", err)
		return res.AbortWithStatus(http.StatusInternalServerError)
	}

	return res.JSON(http.StatusOK, articles)
}

func (c *ArticlesController) GetArticle(res router.Response, req router.Request) error {
	id, err := strconv.Atoi(req.Param("id"))
	if err != nil {
		c.Log.Info("Can't convert article id to int", "err", err.Error())
		return res.AbortWithStatus(http.StatusInternalServerError)
	}

	article, err := c.ArticleInteractor.FindByID(id)
	if err != nil {
		c.Log.Debug("Can't find article", "id", id)
		return res.AbortWithStatus(http.StatusNotFound)
	}

	return res.JSON(http.StatusOK, article)
}
