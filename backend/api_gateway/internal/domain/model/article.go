package model

import "time"

type Article struct {
	ID             string    `json:"id,omitempty" redis:"id"`
	Title          string    `json:"title" redis:"title"`
	Content        string    `json:"content,omitempty" redis:"content"`
	AuthorName     string    `json:"author_name" redis:"author_name"`
	AuthorID       int64     `json:"author_id,omitempty" redis:"author_id"`
	AuthorPhotoURL string    `json:"author_photo_url,omitempty" redis:"author_photo_url"`
	CreatedAt      time.Time `json:"created_at" redis:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" redis:"updated_at"`
}
