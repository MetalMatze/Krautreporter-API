package marshaller

import (
	"time"

	"github.com/metalmatze/krautreporter-api/krautreporter/entity"
)

const KrautreporterURL string = "https://krautreporter.de"

type authorMarshaller struct {
	ID          int       `json:"id"`
	Ordering    int       `json:"order"`
	Name        string    `json:"name"`
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	Biography   string    `json:"biography"`
	SocialMedia string    `json:"socialmedia"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Images map[string][]imageMarshaller `json:"images"`
}

func Author(a *entity.Author) map[string]authorMarshaller {
	return map[string]authorMarshaller{
		"data": marshallAuthor(a),
	}
}

func Authors(authors []*entity.Author) map[string][]authorMarshaller {
	var as []authorMarshaller

	for _, a := range authors {
		as = append(as, marshallAuthor(a))
	}

	return map[string][]authorMarshaller{
		"data": as,
	}
}

func marshallAuthor(a *entity.Author) authorMarshaller {
	am := authorMarshaller{
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

	if len(a.Images) > 0 {
		am.Images = Images(a.Images)
	}

	return am
}
