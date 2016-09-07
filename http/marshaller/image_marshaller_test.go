package marshaller

import (
	"encoding/json"
	"testing"

	"github.com/metalmatze/krautreporter-api/entity"
	"github.com/stretchr/testify/assert"
)

func TestMarshallImage(t *testing.T) {
	i := entity.Image{ID: 123, Width: 256, Src: "/foo.jpg"}

	b, err := json.Marshal(marshallImage(i))
	assert.Nil(t, err)
	assert.JSONEq(t, `{"id":123,"width":256,"src":"https://krautreporter.de/foo.jpg"}`, string(b))
}

func TestImage(t *testing.T) {
	i := []entity.Image{}

	b, err := json.Marshal(Images(i))
	assert.Nil(t, err)
	assert.JSONEq(t, `{"data":null}`, string(b))

	i = append(i, entity.Image{ID: 123, Width: 256, Src: "/foo.jpg"})
	b, err = json.Marshal(Images(i))
	assert.Nil(t, err)
	assert.JSONEq(t, `{"data":[{"id":123,"width":256,"src":"https://krautreporter.de/foo.jpg"}]}`, string(b))
}
