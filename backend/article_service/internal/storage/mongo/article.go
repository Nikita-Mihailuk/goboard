package mongo

import (
	"context"
	"errors"
	"github.com/Nikita-Mihailuk/goboard/backend/article_service/internal/domain/dto"
	"github.com/Nikita-Mihailuk/goboard/backend/article_service/internal/domain/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongoDriver "go.mongodb.org/mongo-driver/v2/mongo"
	"time"
)

const (
	collectionName = "articles"
	dbName         = "article_service_goboard"
)

func (s *Storage) collection() *mongoDriver.Collection {
	return s.db.Database(dbName).Collection(collectionName)
}

func (s *Storage) SaveArticle(ctx context.Context, input dto.CreateArticleInput) error {
	now := time.Now()
	article := model.Article{
		ID:             primitive.NewObjectID().Hex(),
		Title:          input.Title,
		Content:        input.Content,
		AuthorName:     input.AuthorName,
		AuthorID:       input.AuthorID,
		AuthorPhotoURL: input.AuthorPhotoURL,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	_, err := s.collection().InsertOne(ctx, article)
	return err
}

func (s *Storage) GetByID(ctx context.Context, id string) (model.Article, error) {
	var article model.Article
	err := s.collection().FindOne(ctx, bson.M{"_id": id}).Decode(&article)
	if errors.Is(err, mongoDriver.ErrNoDocuments) {
		return model.Article{}, ErrArticleNotFound
	}
	return article, err
}

func (s *Storage) GetByAuthorID(ctx context.Context, authorID int64) ([]model.Article, error) {
	cursor, err := s.collection().Find(ctx, bson.M{"author_id": authorID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var articles []model.Article
	for cursor.Next(ctx) {
		var article model.Article
		if err = cursor.Decode(&article); err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}
	return articles, nil
}

func (s *Storage) GetAll(ctx context.Context) ([]model.Article, error) {
	cursor, err := s.collection().Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var articles []model.Article
	for cursor.Next(ctx) {
		var article model.Article
		if err = cursor.Decode(&article); err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}
	return articles, nil
}

func (s *Storage) UpdateArticle(ctx context.Context, article model.Article) error {
	filter := bson.M{"_id": article.ID}
	update := bson.M{
		"$set": bson.M{
			"title":      article.Title,
			"content":    article.Content,
			"updated_at": article.UpdatedAt,
		},
	}

	_, err := s.collection().UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) DeleteArticle(ctx context.Context, id string) error {
	_, err := s.collection().DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) UpdateAuthorArticle(ctx context.Context, message dto.UpdateAuthorMessage) error {
	filter := bson.M{"author_id": message.UserID}
	update := bson.M{
		"$set": bson.M{
			"author_name":      message.UserName,
			"author_photo_url": message.UserPhotoURL,
			"updated_at":       time.Now(),
		},
	}

	_, err := s.collection().UpdateMany(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}
