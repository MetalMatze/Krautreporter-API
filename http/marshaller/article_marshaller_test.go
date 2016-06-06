package marshaller

import (
	"testing"

	"encoding/json"

	"github.com/MetalMatze/Krautreporter-API/krautreporter/entity"
	"github.com/stretchr/testify/assert"
)

func TestMarshallArticle(t *testing.T) {
	a := &entity.Article{
		ID:       123,
		Ordering: 10,
		Title:    "Title",
		Headline: "Headline",
		Preview:  true,
		URL:      "/123--article",
		Excerpt:  "foo",
		Content:  "bar",
		AuthorID: 13,
	}

	b, err := json.Marshal(Article(a))
	assert.Nil(t, err)
	assert.JSONEq(
		t,
		`{"data":{"id":123,"order":10,"title":"Title","headline":"Headline","date":"0001-01-01T00:00:00Z","morgenpost":false,"preview":true,"url":"https://krautreporter.de/123--article","excerpt":"foo","content":"bar","created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","author_id":13,"images":null}}`,
		string(b),
	)

	a.Images = append(a.Images, entity.Image{ID: 123, Width: 256, Src: "/foo.jpg"})

	b, err = json.Marshal(Article(a))
	assert.Nil(t, err)
	assert.JSONEq(
		t,
		`{"data":{"id":123,"order":10,"title":"Title","headline":"Headline","date":"0001-01-01T00:00:00Z","morgenpost":false,"preview":true,"url":"https://krautreporter.de/123--article","excerpt":"foo","content":"bar","created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","author_id":13,"images":{"data":[{"id":123,"width":256,"src":"https://krautreporter.de/foo.jpg"}]}}}`,
		string(b),
	)
}

func TestMarshallArticles(t *testing.T) {
	articles := []*entity.Article{{
		ID:       1,
		Ordering: 1,
	}, {
		ID:       2,
		Ordering: 0,
	}}

	b, err := json.Marshal(Articles(articles))
	assert.Nil(t, err)
	assert.JSONEq(
		t,
		`{"data":[{"id":1,"order":1,"title":"","headline":"","date":"0001-01-01T00:00:00Z","morgenpost":false,"preview":false,"url":"https://krautreporter.de","excerpt":"","content":"","created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","author_id":0,"images":null},{"id":2,"order":0,"title":"","headline":"","date":"0001-01-01T00:00:00Z","morgenpost":false,"preview":false,"url":"https://krautreporter.de","excerpt":"","content":"","created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","author_id":0,"images":null}]}`,
		string(b),
	)
}
