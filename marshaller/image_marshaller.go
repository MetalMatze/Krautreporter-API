package marshaller

import (
	"github.com/metalmatze/krautreporter-api/entity"
)

// Image is a marshalled struct of the entity Image
type Image struct {
	ID    int    `json:"id"`
	Width int    `json:"width"`
	Src   string `json:"src"`
}

func marshallImage(i entity.Image) Image {
	return Image{
		ID:    i.ID,
		Width: i.Width,
		Src:   KrautreporterURL + i.Src,
	}
}

// FromImages turns a slice of Images into a marshalled data structure
func FromImages(images []entity.Image) map[string][]Image {
	im := []Image{}

	for _, i := range images {
		im = append(im, marshallImage(i))
	}

	return map[string][]Image{
		"data": im,
	}
}
