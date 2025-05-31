package grpc

import (
	"context"
	"github.com/Nikita-Mihailuk/goboard/backend/article_service/internal/domain/dto"
	"github.com/Nikita-Mihailuk/goboard/backend/article_service/internal/domain/model"
	articleServicev1 "github.com/Nikita-Mihailuk/protos_goboard/gen/go/article_service"
	"google.golang.org/grpc"
)

type ArticleService interface {
	CreateArticle(ctx context.Context, input dto.CreateArticleInput) error
	GetArticleByID(ctx context.Context, id string) (model.Article, error)
	GetArticlesByAuthorID(ctx context.Context, authorID int64) ([]model.Article, error)
	GetAllArticles(ctx context.Context) ([]model.Article, error)
	UpdateArticleByID(ctx context.Context, input dto.UpdateArticleInput) error
	DeleteArticleByID(ctx context.Context, id string) error
}

type serverGRPC struct {
	articleService ArticleService
	articleServicev1.UnimplementedUserServer
}

func RegisterGRPCServer(grpcServer *grpc.Server, articleService ArticleService) {
	articleServicev1.RegisterUserServer(grpcServer, &serverGRPC{articleService: articleService})
}
