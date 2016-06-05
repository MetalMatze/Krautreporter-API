package entity

import "time"

type Author struct {
	ID          int
	Ordering    int
	Name        string
	Title       string
	URL         string
	Biography   string
	SocialMedia string
	CreatedAt   time.Time
	UpdatedAt   time.Time

	Images []Image `gorm:"polymorphic:Imageable;"`
	Crawl  Crawl   `gorm:"polymorphic:Crawlable;"`
}
