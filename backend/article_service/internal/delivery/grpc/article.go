package grpc

import (
	"context"
	articleServicev1 "github.com/Nikita-Mihailuk/protos_goboard/gen/go/article_service"
)

func (s *serverGRPC) CreateArticle(ctx context.Context, req *articleServicev1.CreateArticleRequest) (*articleServicev1.CreateArticleResponse, error) {
	panic("implement me")
}

func (s *serverGRPC) GetArticleByID(ctx context.Context, req *articleServicev1.GetArticleByIDRequest) (*articleServicev1.GetArticleByIDResponse, error) {
	panic("implement me")
}

func (s *serverGRPC) GetArticlesByAuthorID(ctx context.Context, req *articleServicev1.GetArticlesByAuthorIDRequest) (*articleServicev1.GetArticlesByAuthorIDResponse, error) {
	panic("implement me")
}

func (s *serverGRPC) GetAllArticle(ctx context.Context, req *articleServicev1.GetAllArticleRequest) (*articleServicev1.GetAllArticleResponse, error) {
	panic("implement me")
}

func (s *serverGRPC) UpdateArticle(ctx context.Context, req *articleServicev1.UpdateArticleRequest) (*articleServicev1.UpdateArticleResponse, error) {
	panic("implement me")
}

func (s *serverGRPC) DeleteArticle(ctx context.Context, req *articleServicev1.DeleteArticleRequest) (*articleServicev1.DeleteArticleResponse, error) {
	panic("implement me")
}
