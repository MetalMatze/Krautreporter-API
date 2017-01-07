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

	Images []Image `gorm:"polymorphic:Imageable;"`
	Crawl  Crawl   `gorm:"polymorphic:Crawlable;"`
}

// AddImage adds an image to the author and makes sure that there's only one image for each width
func (a *Author) AddImage(i Image) {
	for key, image := range a.Images {
		if image.Width == i.Width {
			a.Images[key].Src = i.Src
			return
		}
	}

	a.Images = append(a.Images, i)
}
