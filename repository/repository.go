package repository

import (
	"github.com/go-kit/kit/log"
	"github.com/jinzhu/gorm"
	gocache "github.com/patrickmn/go-cache"
)

type Repository struct {
	Cache  *gocache.Cache
	DB     *gorm.DB
	Logger log.Logger
}
