package entity

import "time"

type Crawl struct {
	ID   int       `json:"id"`
	Next time.Time `json:"next"`

	CrawlableID   int    `json:"-"`
	CrawlableType string `json:"-"`
}
