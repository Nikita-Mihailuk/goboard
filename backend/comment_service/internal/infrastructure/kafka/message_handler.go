package kafka

import (
	"context"
	"encoding/json"
	"github.com/Nikita-Mihailuk/goboard/backend/comment_service/internal/domain/dto"
)

type Handler struct {
	commentUpdater CommentUpdater
}

type CommentUpdater interface {
	UpdateAuthorComments(ctx context.Context, message dto.UpdateAuthorMessage) error
}

func NewHandler(commentUpdater CommentUpdater) *Handler {
	return &Handler{
		commentUpdater: commentUpdater,
	}
}

func (h *Handler) HandleMessage(ctx context.Context, msg []byte) error {
	var authorUpdateMessage dto.UpdateAuthorMessage
	if err := json.Unmarshal(msg, &authorUpdateMessage); err != nil {
		return err
	}

	if err := h.commentUpdater.UpdateAuthorComments(ctx, authorUpdateMessage); err != nil {
		return err
	}

	return nil
}
