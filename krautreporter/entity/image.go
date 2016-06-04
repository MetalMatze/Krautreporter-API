package entity

type Image struct {
	ID    int    `json:"id"`
	Width int    `json:"width"`
	Src   string `json:"src"`

	ImageableID   int    `json:"-"`
	ImageableType string `json:"-"`
}
