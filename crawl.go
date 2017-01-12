package krautreporter

import (
	"math/rand"
	"time"
)

// Crawl is a polymorphic relationship to entities which are being crawled in the future
type Crawl struct {
	ID   int       `json:"id"`
	Next time.Time `json:"next"`

	CrawlableID   int    `json:"-"`
	CrawlableType string `json:"-"`
}

// NextRandom adds to crawl's next 5 hours + random 0 to 30 min
func (c *Crawl) NextRandom() {
	constant := 5 * time.Hour
	variable := 30 * time.Minute
	random := rand.Intn(int(variable.Seconds()))

	dur := constant + time.Duration(random)*time.Second

	c.Next = time.Now().Add(dur)
}
