package marshaller

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/MetalMatze/Krautreporter-API/krautreporter/entity"
	"github.com/stretchr/testify/assert"
)

func TestMarshallAuthor(t *testing.T) {
	a := &entity.Author{
		ID:          13,
		Ordering:    65,
		Name:        "Tilo Jung",
		Title:       "Politik",
		URL:         "/13--tilo-jung",
		Biography:   "Tilo Jung",
		SocialMedia: "TWITTER | FACEBOOK",
		CreatedAt:   time.Date(2016, 06, 05, 22, 32, 0, 0, time.UTC),
		UpdatedAt:   time.Date(2016, 06, 05, 22, 32, 0, 0, time.UTC),
	}

	b, err := json.Marshal(Author(a))
	assert.Nil(t, err)
	assert.JSONEq(
		t,
		`{"data":{"id":13,"order":65,"name":"Tilo Jung","title":"Politik","url":"https://krautreporter.de/13--tilo-jung","biography":"Tilo Jung","socialmedia":"TWITTER | FACEBOOK","created_at":"2016-06-05T22:32:00Z","updated_at":"2016-06-05T22:32:00Z","images":null}}`,
		string(b),
	)

	a.Images = append(a.Images, entity.Image{ID: 123, Width: 256, Src: "/foo.jpg"})

	b, err = json.Marshal(Author(a))
	assert.Nil(t, err)
	assert.JSONEq(
		t,
		`{"data":{"id":13,"order":65,"name":"Tilo Jung","title":"Politik","url":"https://krautreporter.de/13--tilo-jung","biography":"Tilo Jung","socialmedia":"TWITTER | FACEBOOK","created_at":"2016-06-05T22:32:00Z","updated_at":"2016-06-05T22:32:00Z","images":{"data":[{"id":123,"width":256,"src":"https://krautreporter.de/foo.jpg"}]}}}`,
		string(b),
	)
}
