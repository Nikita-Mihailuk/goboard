package comment

import (
	"context"
	"github.com/Nikita-Mihailuk/goboard/backend/comment_service/internal/domain/dto"
	"github.com/Nikita-Mihailuk/goboard/backend/comment_service/internal/domain/model"
)

func (s *CommentService) CreateComment(ctx context.Context, input dto.CreateCommentInput) error {
	panic("implement me")
}

func (s *CommentService) GetCommentsByArticleID(ctx context.Context, articleID string) ([]model.Comment, error) {
	panic("implement me")
}

func (s *CommentService) UpdateCommentByID(ctx context.Context, input dto.UpdateCommentInput) error {
	panic("implement me")
}

func (s *CommentService) DeleteCommentByID(ctx context.Context, id string) error {
	panic("implement me")
}
