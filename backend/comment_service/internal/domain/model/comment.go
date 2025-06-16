package model

import "time"

type Comment struct {
	ID             string     `bson:"_id,omitempty" json:"id"`
	Content        string     `bson:"content" json:"content"`
	ArticleID      string     `bson:"article_id" json:"article_id"`
	AuthorName     string     `bson:"author_name" json:"author_name"`
	AuthorID       int64      `bson:"author_id" json:"author_id"`
	AuthorPhotoURL string     `bson:"author_photo_url,omitempty" json:"author_photo_url,omitempty"`
	ParentID       string     `bson:"parent_id,omitempty" json:"parent_id,omitempty"`
	CreatedAt      time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt      *time.Time `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}
