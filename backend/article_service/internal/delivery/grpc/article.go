package grpc

import (
	"context"
	"errors"
	"github.com/Nikita-Mihailuk/goboard/backend/article_service/internal/domain/dto"
	"github.com/Nikita-Mihailuk/goboard/backend/article_service/internal/service/article"
	articleServicev1 "github.com/Nikita-Mihailuk/protos_goboard/gen/go/article_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *serverGRPC) CreateArticle(ctx context.Context, req *articleServicev1.CreateArticleRequest) (*articleServicev1.CreateArticleResponse, error) {
	inputArticle, err := validateCreateArticleRequest(req)
	if err != nil {
		return nil, err
	}

	err = s.articleService.CreateArticle(ctx, inputArticle)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &articleServicev1.CreateArticleResponse{}, nil
}

func (s *serverGRPC) GetArticleByID(ctx context.Context, req *articleServicev1.GetArticleByIDRequest) (*articleServicev1.GetArticleByIDResponse, error) {
	if req.ArticleId == "" {
		return nil, status.Error(codes.InvalidArgument, "missing fields")
	}

	outputArticle, err := s.articleService.GetArticleByID(ctx, req.GetArticleId())
	if err != nil {
		if errors.Is(err, article.ErrArticleNotFound) {
			return nil, status.Error(codes.NotFound, "article not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &articleServicev1.GetArticleByIDResponse{
		Title:          outputArticle.Title,
		AuthorName:     outputArticle.AuthorName,
		AuthorPhotoUrl: outputArticle.AuthorPhotoURL,
		Content:        outputArticle.Content,
		CreatedAt:      timestamppb.New(outputArticle.CreatedAt),
		UpdatedAt:      timestamppb.New(outputArticle.UpdatedAt),
	}, nil
}

func (s *serverGRPC) GetArticlesByAuthorID(ctx context.Context, req *articleServicev1.GetArticlesByAuthorIDRequest) (*articleServicev1.GetArticlesByAuthorIDResponse, error) {
	if req.GetAuthorId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "missing fields")
	}

	articles, err := s.articleService.GetArticlesByAuthorID(ctx, req.GetAuthorId())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	pbArticles := make([]*articleServicev1.ArticleForList, 0, len(articles))
	for _, artcl := range articles {
		pbArticle := &articleServicev1.ArticleForList{
			Title:      artcl.Title,
			AuthorName: artcl.AuthorName,
			ArticleId:  artcl.ID,
			CreatedAt:  timestamppb.New(artcl.CreatedAt),
			UpdatedAt:  timestamppb.New(artcl.UpdatedAt),
		}
		pbArticles = append(pbArticles, pbArticle)
	}
	return &articleServicev1.GetArticlesByAuthorIDResponse{Articles: pbArticles}, nil
}

func (s *serverGRPC) GetAllArticle(ctx context.Context, req *articleServicev1.GetAllArticleRequest) (*articleServicev1.GetAllArticleResponse, error) {
	articles, err := s.articleService.GetAllArticles(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	pbArticles := make([]*articleServicev1.ArticleForList, 0, len(articles))
	for _, artcl := range articles {
		pbArticle := &articleServicev1.ArticleForList{
			Title:      artcl.Title,
			AuthorName: artcl.AuthorName,
			ArticleId:  artcl.ID,
			CreatedAt:  timestamppb.New(artcl.CreatedAt),
			UpdatedAt:  timestamppb.New(artcl.UpdatedAt),
		}
		pbArticles = append(pbArticles, pbArticle)
	}
	return &articleServicev1.GetAllArticleResponse{Articles: pbArticles}, nil
}

func (s *serverGRPC) UpdateArticle(ctx context.Context, req *articleServicev1.UpdateArticleRequest) (*articleServicev1.UpdateArticleResponse, error) {
	if req.GetArticleId() == "" {
		return nil, status.Error(codes.InvalidArgument, "missing fields")
	}

	inputArticle := dto.UpdateArticleInput{
		Title:   req.GetTitle(),
		Content: req.GetContent(),
		ID:      req.GetArticleId(),
	}

	err := s.articleService.UpdateArticleByID(ctx, inputArticle)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &articleServicev1.UpdateArticleResponse{}, nil
}

func (s *serverGRPC) DeleteArticle(ctx context.Context, req *articleServicev1.DeleteArticleRequest) (*articleServicev1.DeleteArticleResponse, error) {
	if req.GetArticleId() == "" {
		return nil, status.Error(codes.InvalidArgument, "missing fields")
	}

	err := s.articleService.DeleteArticleByID(ctx, req.GetArticleId())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &articleServicev1.DeleteArticleResponse{}, nil
}

func validateCreateArticleRequest(req *articleServicev1.CreateArticleRequest) (dto.CreateArticleInput, error) {
	if req.GetAuthorId() == 0 ||
		req.GetAuthorName() == "" ||
		req.GetTitle() == "" ||
		req.GetContent() == "" {

		return dto.CreateArticleInput{}, status.Errorf(codes.InvalidArgument, "missing fields")
	}

	return dto.CreateArticleInput{
		AuthorID:       req.GetAuthorId(),
		AuthorName:     req.GetAuthorName(),
		Title:          req.GetTitle(),
		Content:        req.GetContent(),
		AuthorPhotoURL: req.GetAuthorPhotoUrl(),
	}, nil
}
