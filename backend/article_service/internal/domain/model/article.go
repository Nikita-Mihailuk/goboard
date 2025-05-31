package model

import "time"

type Article struct {
	ID         string
	Title      string
	Content    string
	AuthorName string
	AuthorID   int64
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
