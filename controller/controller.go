package controller

import (
	"github.com/go-kit/kit/log"
	"github.com/metalmatze/krautreporter-api/repository"
)

type Controller struct {
	Logger     log.Logger
	Repository repository.Repository
}
