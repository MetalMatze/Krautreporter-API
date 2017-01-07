package marshaller

import (
	"time"

	krautreporter "github.com/metalmatze/krautreporter-api"
)

// Article is a marshalled struct of the entity Article
type Article struct {
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

	AuthorID int                `json:"author_id"`
	Images   map[string][]Image `json:"images"`
}

// FromArticle turns a single FromArticle into a marshalled data structure
func FromArticle(a *krautreporter.Article) map[string]Article {
	return map[string]Article{
		"data": marshallArticle(a),
	}
}

// FromArticles turns a slice of FromArticles into a marshalled data structure
func FromArticles(articles []*krautreporter.Article) map[string][]Article {
	var as []Article

	for _, a := range articles {
		as = append(as, marshallArticle(a))
	}

	return map[string][]Article{
		"data": as,
	}
}

func marshallArticle(a *krautreporter.Article) Article {
	am := Article{
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

	am.Images = FromImages(a.Images)

	return am
}
