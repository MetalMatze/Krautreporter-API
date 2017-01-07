package krautreporter

import "time"

// Crawl is a polymorphic relationship to entities which are being crawled in the future
type Crawl struct {
	ID   int       `json:"id"`
	Next time.Time `json:"next"`

	CrawlableID   int    `json:"-"`
	CrawlableType string `json:"-"`
}
