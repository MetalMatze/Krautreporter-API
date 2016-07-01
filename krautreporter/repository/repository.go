package repository

import (
	"github.com/go-kit/kit/log"
	"github.com/jinzhu/gorm"
	gocache "github.com/patrickmn/go-cache"
)

type repository struct {
	cache  *gocache.Cache
	db     *gorm.DB
	logger log.Logger
}

func newRepository(logger log.Logger, db *gorm.DB, cache *gocache.Cache) repository {
	return repository{logger: logger, db: db, cache: cache}
}
