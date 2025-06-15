package model

import "time"

type Comment struct {
	ID             string     `bson:"_id,omitempty"`
	Content        string     `bson:"content"`
	ArticleID      string     `bson:"article_id"`
	AuthorName     string     `bson:"author_name"`
	AuthorID       int64      `bson:"author_id"`
	AuthorPhotoURL string     `bson:"author_photo_url,omitempty"`
	ParentID       string     `bson:"parent_id,omitempty"`
	CreatedAt      time.Time  `bson:"created_at"`
	UpdatedAt      *time.Time `bson:"updated_at,omitempty"`
}
