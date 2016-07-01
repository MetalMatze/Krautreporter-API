package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CrawlsController struct {
	*Controller
}

func (c *CrawlsController) GetCrawls(ctx *gin.Context) {
	crawls, err := c.interactor.NextCrawls(20)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}

	ctx.JSON(http.StatusOK, crawls)
}
