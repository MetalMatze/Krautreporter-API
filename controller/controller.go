package controller

import (
	"github.com/go-kit/kit/log"
	"github.com/metalmatze/krautreporter-api/repository"
)

// Controller has a logger and talks to the Repository to retrieve data from the database
type Controller struct {
	Logger     log.Logger
	Repository repository.Repository
}
