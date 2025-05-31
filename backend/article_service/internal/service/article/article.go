package article

import (
	"context"
	"errors"
	"github.com/Nikita-Mihailuk/goboard/backend/article_service/internal/domain/dto"
	"github.com/Nikita-Mihailuk/goboard/backend/article_service/internal/domain/model"
	"github.com/Nikita-Mihailuk/goboard/backend/article_service/internal/storage/mongo"
	"go.uber.org/zap"
	"time"
)

func (s *ArticleService) CreateArticle(ctx context.Context, input dto.CreateArticleInput) error {
	err := s.articleSaver.SaveArticle(ctx, input)
	if err != nil {
		s.log.Error("failed save article", zap.Error(err))
		return err
	}
	s.log.Info("save article", zap.String("title", input.Title))
	return nil
}

func (s *ArticleService) GetArticleByID(ctx context.Context, id string) (model.Article, error) {
	article, err := s.articleProvider.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrArticleNotFound) {
			s.log.Error("article not found", zap.String("id", id))
			return model.Article{}, ErrArticleNotFound
		}
		s.log.Error("failed get article", zap.Error(err))
	}
	return article, err
}

func (s *ArticleService) GetArticlesByAuthorID(ctx context.Context, authorID int64) ([]model.Article, error) {
	articles, err := s.articleProvider.GetByAuthorID(ctx, authorID)
	if err != nil {
		s.log.Error("failed get articles by author id", zap.Error(err))
		return nil, err
	}
	s.log.Info("get articles by author id", zap.Int64("author_id", authorID))
	return articles, nil
}

func (s *ArticleService) GetAllArticles(ctx context.Context) ([]model.Article, error) {
	articles, err := s.articleProvider.GetAll(ctx)
	if err != nil {
		s.log.Error("failed get all articles", zap.Error(err))
		return nil, err
	}
	return articles, nil
}

func (s *ArticleService) UpdateArticleByID(ctx context.Context, input dto.UpdateArticleInput) error {
	article, err := s.articleProvider.GetByID(ctx, input.ID)
	if err != nil {
		s.log.Error("failed get article", zap.Error(err))
		return err
	}

	if input.Title != "" {
		article.Title = input.Title
	}

	if input.Content != "" {
		article.Content = input.Content
	}

	article.UpdatedAt = time.Now()

	err = s.articleUpdater.UpdateArticle(ctx, article)
	if err != nil {
		s.log.Error("failed update article", zap.Error(err))
		return err
	}
	s.log.Info("update article", zap.String("id", input.ID))
	return nil
}

func (s *ArticleService) DeleteArticleByID(ctx context.Context, id string) error {
	err := s.articleDeleter.DeleteArticle(ctx, id)
	if err != nil {
		s.log.Error("failed delete article", zap.Error(err))
		return err
	}
	s.log.Info("delete article", zap.String("id", id))
	return nil
}
