package entity

import "time"

type Article struct {
	ID        int
	Ordering  int
	Title     string
	Headline  string
	Date      string
	Preview   bool
	URL       string
	Excerpt   string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time

	AuthorID int
	Author   *Author
	Images   []Image `gorm:"polymorphic:Imageable;"`
	Crawl    Crawl   `gorm:"polymorphic:Crawlable;"`
}
