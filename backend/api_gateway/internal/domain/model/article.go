package model

import "time"

type Article struct {
	ID             string    `json:"id,omitempty"`
	Title          string    `json:"title"`
	Content        string    `json:"content,omitempty"`
	AuthorName     string    `json:"author_name"`
	AuthorID       int64     `json:"author_id,omitempty"`
	AuthorPhotoURL string    `json:"author_photo_url,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
