package mongo

import (
	"context"
	"github.com/Nikita-Mihailuk/goboard/backend/comment_service/internal/domain/dto"
	"github.com/Nikita-Mihailuk/goboard/backend/comment_service/internal/domain/model"
	mongoDriver "go.mongodb.org/mongo-driver/v2/mongo"
)

const (
	collectionName = "comments"
	dbName         = "goboard_db"
)

func (s *Storage) collection() *mongoDriver.Collection {
	return s.db.Database(dbName).Collection(collectionName)
}

func (s *Storage) SaveComment(ctx context.Context, input dto.CreateCommentInput) error {
	panic("implement me")
}

func (s *Storage) GetByArticleID(ctx context.Context, articleID string) ([]model.Comment, error) {
	panic("implement me")
}

func (s *Storage) UpdateComment(ctx context.Context, comment dto.UpdateCommentInput) error {
	panic("implement me")
}

func (s *Storage) DeleteComment(ctx context.Context, id string) error {
	panic("implement me")
}

func (s *Storage) UpdateAuthorComments(ctx context.Context, message dto.UpdateAuthorMessage) error {
	panic("implement me")
}
