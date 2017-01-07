package entity

// Image is polymorphic and can be used by Article or Author
type Image struct {
	ID    int
	Width int
	Src   string

	ImageableID   int
	ImageableType string
}
