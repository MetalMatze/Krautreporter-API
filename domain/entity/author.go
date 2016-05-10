package entity

import "time"

type AuthorRepository interface {
	Find() ([]*Author, error)
	FindByID(int) (*Author, error)
	SaveAll([]Author) error
}

type Author struct {
	ID          int       `gorm:"primary_key" json:"id"`
	Ordering    int       `json:"order"`
	Name        string    `json:"name"`
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	Biography   string    `json:"biography"`
	SocialMedia string    `json:"social-media"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
