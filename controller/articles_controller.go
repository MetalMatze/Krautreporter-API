package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/metalmatze/krautreporter-api/marshaller"
	"github.com/metalmatze/krautreporter-api/repository"
)

// ArticlesPerPage is the number of articles to be on a page with pagination
const ArticlesPerPage int = 10

// GetArticles returns a paginated list of marshalled FromArticles
func (ctrl *Controller) GetArticles(c *gin.Context) {
	id := repository.MaxArticleID
	if c.Query("olderthan") != "" {
		olderthan, err := strconv.Atoi(c.Query("olderthan"))
		if err != nil {
			ctrl.Logger.Log("msg", "Can't convert olderthan id to int", "err", err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		id = olderthan
	}

	articles, err := ctrl.Repository.FindArticlesOlderThan(id, ArticlesPerPage)
	if err != nil {
		if err == repository.ErrArticleNotFound {
			ctrl.Logger.Log("msg", "Can't find olderthan article", "id", id)
			status := http.StatusNotFound
			c.JSON(status, gin.H{
				"message":     http.StatusText(status),
				"status_code": status,
			})
		}

		ctrl.Logger.Log("msg", "Failed to get olderthan articles", "id", id, "err", err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	c.JSON(http.StatusOK, marshaller.FromArticles(articles))
}

// GetArticle returns a single marshalled article
func (ctrl *Controller) GetArticle(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ctrl.Logger.Log("msg", "Can't convert article id to int", "err", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	article, err := ctrl.Repository.FindArticleByID(id)
	if err != nil {
		ctrl.Logger.Log("msg", "Can't find article", "id", id)
		c.AbortWithStatus(http.StatusNotFound)
	}

	c.JSON(http.StatusOK, marshaller.FromArticle(article))
}
