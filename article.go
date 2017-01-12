package krautreporter

import "time"

// Article posted on krautreporter.de
type Article struct {
	ID        int
	Ordering  int
	Title     string
	Headline  string
	Date      time.Time
	Preview   bool
	URL       string
	Excerpt   string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time

	AuthorID int
	Author   *Author
	Images   []*Image `gorm:"polymorphic:Imageable;"`
	Crawl    *Crawl   `gorm:"polymorphic:Crawlable;"`
}

// AddImage adds an image to the article and makes sure that there's only one image for each width
func (a *Article) AddImage(i *Image) {
	for key, image := range a.Images {
		if image.Width == i.Width {
			a.Images[key].Src = i.Src
			return
		}
	}

	a.Images = append(a.Images, i)
}
