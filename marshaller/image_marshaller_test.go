package marshaller

import (
	"encoding/json"
	"testing"

	krautreporter "github.com/metalmatze/krautreporter-api"
	"github.com/stretchr/testify/assert"
)

func TestMarshallImage(t *testing.T) {
	i := krautreporter.Image{ID: 123, Width: 256, Src: "/foo.jpg"}

	b, err := json.Marshal(marshallImage(i))
	assert.Nil(t, err)
	assert.JSONEq(t, `{"id":123,"width":256,"src":"https://krautreporter.de/foo.jpg"}`, string(b))
}

func TestImage(t *testing.T) {
	i := []krautreporter.Image{}

	b, err := json.Marshal(FromImages(i))
	assert.Nil(t, err)
	assert.JSONEq(t, `{"data":[]}`, string(b))

	i = append(i, krautreporter.Image{ID: 123, Width: 256, Src: "/foo.jpg"})
	b, err = json.Marshal(FromImages(i))
	assert.Nil(t, err)
	assert.JSONEq(t, `{"data":[{"id":123,"width":256,"src":"https://krautreporter.de/foo.jpg"}]}`, string(b))
}
