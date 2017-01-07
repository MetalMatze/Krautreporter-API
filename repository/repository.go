package repository

import (
	"github.com/go-kit/kit/log"
	"github.com/jinzhu/gorm"
	gocache "github.com/patrickmn/go-cache"
)

// Repository reads from the database or the cache if data is present
type Repository struct {
	Cache  *gocache.Cache
	DB     *gorm.DB
	Logger log.Logger
}
