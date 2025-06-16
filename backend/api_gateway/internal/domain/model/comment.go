package model

import "time"

type Comment struct {
	ID             string     `json:"id"`
	Content        string     `json:"content"`
	ArticleID      string     `json:"article_id"`
	AuthorName     string     `json:"author_name"`
	AuthorID       int64      `json:"author_id"`
	AuthorPhotoURL string     `json:"author_photo_url,omitempty"`
	ParentID       string     `json:"parent_id,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
}
