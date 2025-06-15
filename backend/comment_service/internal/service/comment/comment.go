package comment

import (
	"context"
	"github.com/Nikita-Mihailuk/goboard/backend/comment_service/internal/domain/dto"
	"github.com/Nikita-Mihailuk/goboard/backend/comment_service/internal/domain/model"
	"go.uber.org/zap"
)

func (s *CommentService) CreateComment(ctx context.Context, input dto.CreateCommentInput) error {
	if err := s.commentSaver.SaveComment(ctx, input); err != nil {
		s.log.Error("failed save comment", zap.Error(err))
		return err
	}

	return nil
}

func (s *CommentService) GetCommentsByArticleID(ctx context.Context, articleID string) ([]model.Comment, error) {
	comments, err := s.commentProvider.GetByArticleID(ctx, articleID)
	if err != nil {
		s.log.Error("failed get comments", zap.Error(err))
		return nil, err
	}

	return comments, nil
}

func (s *CommentService) UpdateCommentByID(ctx context.Context, input dto.UpdateCommentInput) error {
	if err := s.commentUpdater.UpdateComment(ctx, input); err != nil {
		s.log.Error("failed update comment", zap.Error(err))
		return err
	}

	return nil
}

func (s *CommentService) DeleteCommentByID(ctx context.Context, id string) error {
	if err := s.commentDeleter.DeleteComment(ctx, id); err != nil {
		s.log.Error("failed delete comment", zap.Error(err))
		return err
	}

	return nil
}
