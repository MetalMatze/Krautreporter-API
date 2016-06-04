package entity

import "time"

type Article struct {
	ID        int       `json:"id"`
	Ordering  int       `json:"order"`
	Title     string    `json:"title"`
	Headline  string    `json:"headline"`
	Date      string    `json:"date"`
	Preview   bool      `json:"preview"`
	URL       string    `json:"url"`
	Excerpt   string    `json:"excerpt"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	AuthorID int     `json:"author_id"`
	Author   *Author `json:"-"`
	Images   []Image `gorm:"polymorphic:Imageable;" json:"-"`
	Crawl    Crawl   `gorm:"polymorphic:Crawlable;" json:"-"`
}
