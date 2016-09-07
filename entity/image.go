package entity

type Image struct {
	ID    int
	Width int
	Src   string

	ImageableID   int
	ImageableType string
}
