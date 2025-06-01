package model

import "time"

type Article struct {
	ID             string    `bson:"_id,omitempty"`
	Title          string    `bson:"title"`
	Content        string    `bson:"content"`
	AuthorName     string    `bson:"author_name"`
	AuthorID       int64     `bson:"author_id"`
	AuthorPhotoURL string    `bson:"author_photo_url"`
	CreatedAt      time.Time `bson:"created_at"`
	UpdatedAt      time.Time `bson:"updated_at"`
}
