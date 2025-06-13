package dto

type CreateArticleInput struct {
	Title          string
	AuthorName     string
	AuthorPhotoURL string
	AuthorID       int64
	Content        string
}

type UpdateArticleInput struct {
	ID      string
	Title   string
	Content string
}

type UpdateAuthorMessage struct {
	UserID       int64  `json:"user_id"`
	UserName     string `json:"user_name"`
	UserPhotoURL string `json:"user_photo_url"`
}
