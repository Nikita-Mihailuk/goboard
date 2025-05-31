package article

import (
	"context"
	"github.com/Nikita-Mihailuk/goboard/backend/article_service/internal/domain/dto"
	"github.com/Nikita-Mihailuk/goboard/backend/article_service/internal/domain/model"
)

func (s *ArticleService) CreateArticle(ctx context.Context, input dto.CreateArticleInput) error {
	panic("implement me")
}

func (s *ArticleService) GetArticleByID(ctx context.Context, id string) (model.Article, error) {
	panic("implement me")
}

func (s *ArticleService) GetArticlesByAuthorID(ctx context.Context, authorID int64) ([]model.Article, error) {
	panic("implement me")
}

func (s *ArticleService) GetAllArticles(ctx context.Context) ([]model.Article, error) {
	panic("implement me")
}

func (s *ArticleService) UpdateArticleByID(ctx context.Context, input dto.UpdateArticleInput) error {
	panic("implement me")
}

func (s *ArticleService) DeleteArticleByID(ctx context.Context, id string) error {
	panic("implement me")
}
