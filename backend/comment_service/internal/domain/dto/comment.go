package dto

type CreateCommentInput struct {
	Content        string `json:"content"`
	ArticleID      string `json:"article_id"`
	AuthorName     string `json:"author_name"`
	AuthorID       int64  `json:"author_id"`
	AuthorPhotoURL string `json:"author_photo_url,omitempty"`
	ParentID       string `json:"parent_id,omitempty"`
}

type UpdateCommentInput struct {
	ID      string `json:"id,omitempty"`
	Content string `json:"content"`
}

type UpdateAuthorMessage struct {
	UserID       int64  `json:"user_id"`
	UserName     string `json:"user_name"`
	UserPhotoURL string `json:"user_photo_url"`
}
