package mongo

import (
	"context"
	"github.com/Nikita-Mihailuk/goboard/backend/comment_service/internal/domain/dto"
	"github.com/Nikita-Mihailuk/goboard/backend/comment_service/internal/domain/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongoDriver "go.mongodb.org/mongo-driver/v2/mongo"
	"time"
)

const (
	collectionName = "comments"
	dbName         = "goboard_db"
)

func (s *Storage) collection() *mongoDriver.Collection {
	return s.db.Database(dbName).Collection(collectionName)
}

func (s *Storage) SaveComment(ctx context.Context, input dto.CreateCommentInput) error {
	now := time.Now()
	comment := &model.Comment{
		ID:             primitive.NewObjectID().Hex(),
		Content:        input.Content,
		ArticleID:      input.ArticleID,
		AuthorName:     input.AuthorName,
		AuthorID:       input.AuthorID,
		AuthorPhotoURL: input.AuthorPhotoURL,
		ParentID:       input.ParentID,
		CreatedAt:      now,
	}

	_, err := s.collection().InsertOne(ctx, comment)
	return err
}

func (s *Storage) GetByArticleID(ctx context.Context, articleID string) ([]model.Comment, error) {
	cursor, err := s.collection().Find(ctx, bson.M{"article_id": articleID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var comments []model.Comment
	if err = cursor.All(ctx, &comments); err != nil {
		return nil, err
	}

	return comments, nil
}

func (s *Storage) UpdateComment(ctx context.Context, comment dto.UpdateCommentInput) error {
	filter := bson.M{"_id": comment.ID}
	update := bson.M{
		"$set": bson.M{
			"content":    comment.Content,
			"updated_at": time.Now(),
		},
	}

	_, err := s.collection().UpdateOne(ctx, filter, update)
	return err
}

func (s *Storage) DeleteComment(ctx context.Context, id string) error {
	_, err := s.collection().DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (s *Storage) UpdateAuthorComments(ctx context.Context, message dto.UpdateAuthorMessage) error {
	filter := bson.M{"author_id": message.UserID}
	update := bson.M{
		"$set": bson.M{
			"author_name":      message.UserName,
			"author_photo_url": message.UserPhotoURL,
		},
	}

	_, err := s.collection().UpdateMany(ctx, filter, update)
	return err
}
