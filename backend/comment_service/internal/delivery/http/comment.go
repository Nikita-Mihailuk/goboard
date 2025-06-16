package http

import (
	"encoding/json"
	"net/http"

	"github.com/Nikita-Mihailuk/goboard/backend/comment_service/internal/domain/dto"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) createComment(w http.ResponseWriter, r *http.Request) {
	var input dto.CreateCommentInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	err := h.commentService.CreateComment(r.Context(), input)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) getCommentsByArticleID(w http.ResponseWriter, r *http.Request) {
	articleID := chi.URLParam(r, "article_id")
	comments, err := h.commentService.GetCommentsByArticleID(r.Context(), articleID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(comments)
}

func (h *Handler) updateComment(w http.ResponseWriter, r *http.Request) {
	commentID := chi.URLParam(r, "id")
	var input dto.UpdateCommentInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	input.ID = commentID
	err := h.commentService.UpdateCommentByID(r.Context(), input)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) deleteComment(w http.ResponseWriter, r *http.Request) {
	commentID := chi.URLParam(r, "id")
	err := h.commentService.DeleteCommentByID(r.Context(), commentID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
