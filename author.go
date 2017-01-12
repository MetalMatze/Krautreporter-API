package krautreporter

import "time"

// Author posted on krautreporter.de
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

	Images []*Image `gorm:"polymorphic:Imageable;"`
	Crawl  *Crawl   `gorm:"polymorphic:Crawlable;"`
}

// AddImage adds an image to the author and makes sure that there's only one image for each width
func (a *Author) AddImage(i *Image) {
	for key, image := range a.Images {
		if image.Width == i.Width {
			a.Images[key].Src = i.Src
			return
		}
	}

	a.Images = append(a.Images, i)
}

// NextCrawl merges the current crawl and a new one's timestamp
// If no crawl exists for the author yet, the passed one is taken
func (a *Author) NextCrawl(c *Crawl) {
	if a.Crawl == nil || a.Crawl.ID == 0 {
		a.Crawl = c
		return
	}

	if a.Crawl.ID != 0 && c.ID == 0 {
		a.Crawl.Next = c.Next
		return
	}

	a.Crawl = c
}
