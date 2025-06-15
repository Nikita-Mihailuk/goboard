package comment

import (
	"context"

	"github.com/Nikita-Mihailuk/goboard/backend/comment_service/internal/domain/dto"
	"github.com/Nikita-Mihailuk/goboard/backend/comment_service/internal/domain/model"
	"go.uber.org/zap"
)

type CommentService struct {
	log             *zap.Logger
	commentSaver    CommentSaver
	commentProvider CommentProvider
	commentUpdater  CommentUpdater
	commentDeleter  CommentDeleter
}

func NewCommentService(
	log *zap.Logger,
	commentSaver CommentSaver,
	commentProvider CommentProvider,
	commentUpdater CommentUpdater,
	commentDeleter CommentDeleter,
) *CommentService {
	return &CommentService{
		log:             log,
		commentSaver:    commentSaver,
		commentProvider: commentProvider,
		commentUpdater:  commentUpdater,
		commentDeleter:  commentDeleter,
	}
}

type CommentSaver interface {
	SaveComment(ctx context.Context, input dto.CreateCommentInput) error
}

type CommentProvider interface {
	GetByArticleID(ctx context.Context, articleID string) ([]model.Comment, error)
}

type CommentUpdater interface {
	UpdateComment(ctx context.Context, comment dto.UpdateCommentInput) error
}

type CommentDeleter interface {
	DeleteComment(ctx context.Context, id string) error
}
