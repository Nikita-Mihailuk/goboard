package dto

type CreateArticleInput struct {
	Title          string `json:"title"`
	AuthorName     string `json:"author_name"`
	AuthorPhotoURL string `json:"author_photo_url,omitempty"`
	AuthorID       int64  `json:"author_id"`
	Content        string `json:"content"`
}

type UpdateArticleInput struct {
	ID      string `json:"id"`
	Title   string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
}
