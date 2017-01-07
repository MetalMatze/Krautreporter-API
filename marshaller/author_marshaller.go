package marshaller

import (
	"time"

	krautreporter "github.com/metalmatze/krautreporter-api"
)

// KrautreporterURL is used as a prefix for URLs being marshalled
const KrautreporterURL string = "https://krautreporter.de"

// Author is a marshalled struct of the entity Author
type Author struct {
	ID          int       `json:"id"`
	Ordering    int       `json:"order"`
	Name        string    `json:"name"`
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	Biography   string    `json:"biography"`
	SocialMedia string    `json:"socialmedia"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Images map[string][]Image `json:"images"`
}

// FromAuthor turns a single Author into a marshalled data structure
func FromAuthor(a *krautreporter.Author) map[string]Author {
	return map[string]Author{
		"data": marshallAuthor(a),
	}
}

// FromAuthors turns a slice of Authors into a marshalled data structure
func FromAuthors(authors []*krautreporter.Author) map[string][]Author {
	var as []Author

	for _, a := range authors {
		as = append(as, marshallAuthor(a))
	}

	return map[string][]Author{
		"data": as,
	}
}

func marshallAuthor(a *krautreporter.Author) Author {
	am := Author{
		ID:          a.ID,
		Ordering:    a.Ordering,
		Name:        a.Name,
		Title:       a.Title,
		URL:         KrautreporterURL + a.URL,
		Biography:   a.Biography,
		SocialMedia: a.SocialMedia,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
	}

	am.Images = FromImages(a.Images)

	return am
}
