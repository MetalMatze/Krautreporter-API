package controller

import (
	"net/http"
	"strconv"

	"github.com/MetalMatze/Krautreporter-API/http/marshaller"
	"github.com/MetalMatze/Krautreporter-API/krautreporter/repository"
	"github.com/gin-gonic/gin"
)

const ArticlesPerPage int = 10

type ArticlesController struct {
	*Controller
}

func (c *ArticlesController) GetArticles(ctx *gin.Context) {
	id := repository.MaxArticleID
	if ctx.Query("olderthan") != "" {
		olderthan, err := strconv.Atoi(ctx.Query("olderthan"))
		if err != nil {
			c.logger.Log("msg", "Can't convert olderthan id to int", "err", err.Error())
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
		id = olderthan
	}

	articles, err := c.interactor.ArticlesOlderThan(id, ArticlesPerPage)
	if err != nil {
		if err == repository.ErrArticleNotFound {
			c.logger.Log("msg", "Can't find olderthan article", "id", id)
			status := http.StatusNotFound
			ctx.JSON(status, gin.H{
				"message":     http.StatusText(status),
				"status_code": status,
			})
		}

		c.logger.Log("msg", "Failed to get olderthan articles", "id", id, "err", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}

	ctx.JSON(http.StatusOK, marshaller.Articles(articles))
}

func (c *ArticlesController) GetArticle(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		c.logger.Log("msg", "Can't convert article id to int", "err", err.Error())
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}

	article, err := c.interactor.ArticleByID(id)
	if err != nil {
		c.logger.Log("msg", "Can't find article", "id", id)
		ctx.AbortWithStatus(http.StatusNotFound)
	}

	ctx.JSON(http.StatusOK, marshaller.Article(article))
}
