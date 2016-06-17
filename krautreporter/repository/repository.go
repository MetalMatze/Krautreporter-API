package repository

import (
	"github.com/gollection/gollection/cache"
	"github.com/gollection/gollection/log"
	"github.com/jinzhu/gorm"
)

type repository struct {
	cache cache.Cache
	db    *gorm.DB
	log   log.Logger
}

func newRepository(c cache.Cache, db *gorm.DB, log log.Logger) repository {
	return repository{
		cache: c,
		db:    db,
		log:   log,
	}
}
