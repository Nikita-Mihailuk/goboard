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
