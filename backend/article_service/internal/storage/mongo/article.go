package mongo

import (
	"context"
	"github.com/Nikita-Mihailuk/goboard/backend/article_service/internal/domain/dto"
	"github.com/Nikita-Mihailuk/goboard/backend/article_service/internal/domain/model"
)

func (s *Storage) SaveArticle(ctx context.Context, input dto.CreateArticleInput) error {
	panic("implement me")
}

func (s *Storage) GetByID(ctx context.Context, id string) (model.Article, error) {
	panic("implement me")
}

func (s *Storage) GetByAuthorID(ctx context.Context, authorID int64) ([]model.Article, error) {
	panic("implement me")
}

func (s *Storage) GetAll(ctx context.Context) ([]model.Article, error) {
	panic("implement me")
}

func (s *Storage) UpdateArticle(ctx context.Context, input dto.UpdateArticleInput) error {
	panic("implement me")
}

func (s *Storage) DeleteArticle(ctx context.Context, id string) error {
	panic("implement me")
}
