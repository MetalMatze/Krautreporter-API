package marshaller

import (
	"time"

	"github.com/MetalMatze/Krautreporter-API/krautreporter/entity"
)

type articleMarshaller struct {
	ID         int       `json:"id"`
	Ordering   int       `json:"order"`
	Title      string    `json:"title"`
	Headline   string    `json:"headline"`
	Date       time.Time `json:"date"`
	Morgenpost bool      `json:"morgenpost"`
	Preview    bool      `json:"preview"`
	URL        string    `json:"url"`
	Excerpt    string    `json:"excerpt"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	AuthorID int                          `json:"author_id"`
	Images   map[string][]imageMarshaller `json:"images"`
}

func Article(a *entity.Article) map[string]articleMarshaller {
	return map[string]articleMarshaller{
		"data": marshallArticle(a),
	}
}

func Articles(articles []*entity.Article) map[string][]articleMarshaller {
	var as []articleMarshaller

	for _, a := range articles {
		as = append(as, marshallArticle(a))
	}

	return map[string][]articleMarshaller{
		"data": as,
	}
}

func marshallArticle(a *entity.Article) articleMarshaller {
	am := articleMarshaller{
		ID:         a.ID,
		Ordering:   a.Ordering,
		Title:      a.Title,
		Headline:   a.Headline,
		Date:       a.Date,
		Morgenpost: false,
		Preview:    a.Preview,
		URL:        KrautreporterURL + a.URL,
		Excerpt:    a.Excerpt,
		Content:    a.Content,
		CreatedAt:  a.CreatedAt,
		UpdatedAt:  a.UpdatedAt,

		AuthorID: a.AuthorID,
	}

	if len(a.Images) > 0 {
		am.Images = Images(a.Images)
	}

	return am
}
