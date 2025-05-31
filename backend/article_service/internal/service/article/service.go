package article

import (
	"context"
	"github.com/Nikita-Mihailuk/goboard/backend/article_service/internal/domain/dto"
	"github.com/Nikita-Mihailuk/goboard/backend/article_service/internal/domain/model"
	"go.uber.org/zap"
)

type ArticleService struct {
	log             *zap.Logger
	articleSaver    ArticleSaver
	articleProvider ArticleProvider
	articleUpdater  ArticleUpdater
	articleDeleter  ArticleDeleter
}

func NewArticleService(
	log *zap.Logger,
	articleSaver ArticleSaver,
	articleProvider ArticleProvider,
	articleUpdater ArticleUpdater,
	articleDeleter ArticleDeleter,
) *ArticleService {

	return &ArticleService{
		log:             log,
		articleSaver:    articleSaver,
		articleProvider: articleProvider,
		articleUpdater:  articleUpdater,
		articleDeleter:  articleDeleter,
	}
}

type ArticleSaver interface {
	SaveArticle(ctx context.Context, input dto.CreateArticleInput) error
}

type ArticleProvider interface {
	GetByID(ctx context.Context, id string) (model.Article, error)
	GetByAuthorID(ctx context.Context, authorID int64) ([]model.Article, error)
	GetAll(ctx context.Context) ([]model.Article, error)
}

type ArticleUpdater interface {
	UpdateArticle(ctx context.Context, input dto.UpdateArticleInput) error
}

type ArticleDeleter interface {
	DeleteArticle(ctx context.Context, id string) error
}
