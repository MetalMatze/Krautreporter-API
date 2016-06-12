package controller

import (
	"net/http"

	"github.com/gollection/gollection/router"
)

type CrawlsController struct {
	*Controller
}

func (c *CrawlsController) GetCrawls(res router.Response, req router.Request) error {
	crawls, err := c.interactor.NextCrawls(20)
	if err != nil {
		return res.AbortWithStatus(http.StatusInternalServerError)

	}

	return res.JSON(http.StatusOK, crawls)
}
