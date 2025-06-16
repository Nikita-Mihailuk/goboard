package http

import (
	"context"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/Nikita-Mihailuk/goboard/backend/comment_service/internal/domain/dto"
	"github.com/Nikita-Mihailuk/goboard/backend/comment_service/internal/domain/model"
)

type Handler struct {
	commentService CommentService
}

func NewHandler(commentService CommentService) *Handler {
	return &Handler{commentService: commentService}
}

type CommentService interface {
	CreateComment(ctx context.Context, input dto.CreateCommentInput) error
	GetCommentsByArticleID(ctx context.Context, articleID string) ([]model.Comment, error)
	UpdateCommentByID(ctx context.Context, input dto.UpdateCommentInput) error
	DeleteCommentByID(ctx context.Context, id string) error
}

func (h *Handler) InitRoutes() http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Route("/comments", func(r chi.Router) {
		r.Post("/", h.createComment)
		r.Get("/article/{article_id}", h.getCommentsByArticleID)
		r.Patch("/{id}", h.updateComment)
		r.Delete("/{id}", h.deleteComment)
	})

	return router
}
