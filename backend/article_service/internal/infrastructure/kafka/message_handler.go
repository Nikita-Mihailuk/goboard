package kafka

import (
	"context"
	"encoding/json"
	"github.com/Nikita-Mihailuk/goboard/backend/article_service/internal/domain/dto"
)

type Handler struct {
	articleUpdater ArticleUpdater
}

type ArticleUpdater interface {
	UpdateAuthorArticle(ctx context.Context, message dto.UpdateAuthorMessage) error
}

func NewHandler(articleUpdater ArticleUpdater) *Handler {
	return &Handler{
		articleUpdater: articleUpdater,
	}
}

func (h *Handler) HandleMessage(ctx context.Context, msg []byte) error {
	var authorUpdateMessage dto.UpdateAuthorMessage
	if err := json.Unmarshal(msg, &authorUpdateMessage); err != nil {
		return err
	}

	if err := h.articleUpdater.UpdateAuthorArticle(ctx, authorUpdateMessage); err != nil {
		return err
	}

	return nil
}
