package article_service

import (
	"context"
	"github.com/Nikita-Mihailuk/goboard/backend/api_gateway/internal/domain/dto"
	"github.com/Nikita-Mihailuk/goboard/backend/api_gateway/internal/domain/model"
	articleServicev1 "github.com/Nikita-Mihailuk/protos_goboard/gen/go/article_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ArticleClient struct {
	api articleServicev1.ArticleClient
}

func NewArticleClient(ctx context.Context, addr string) (*ArticleClient, error) {
	cc, err := grpc.DialContext(ctx, addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &ArticleClient{api: articleServicev1.NewArticleClient(cc)}, nil
}

func (c *ArticleClient) CreateArticle(ctx context.Context, input dto.CreateArticleInput) error {
	_, err := c.api.CreateArticle(ctx, &articleServicev1.CreateArticleRequest{
		Title:          input.Title,
		Content:        input.Content,
		AuthorName:     input.AuthorName,
		AuthorId:       input.AuthorID,
		AuthorPhotoUrl: input.AuthorPhotoURL,
	})

	if err != nil {
		_, ok := status.FromError(err)
		if ok {
			return ErrInternalGRPCServer
		}
		return err
	}

	return nil
}

func (c *ArticleClient) GetArticleByID(ctx context.Context, id string) (model.Article, error) {
	resp, err := c.api.GetArticleByID(ctx, &articleServicev1.GetArticleByIDRequest{
		ArticleId: id,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				return model.Article{}, ErrArticleNotFound
			default:
				return model.Article{}, ErrInternalGRPCServer
			}
		}
		return model.Article{}, err
	}

	return model.Article{
		Title:          resp.GetTitle(),
		Content:        resp.GetContent(),
		AuthorName:     resp.GetAuthorName(),
		AuthorPhotoURL: resp.GetAuthorPhotoUrl(),
		AuthorID:       resp.GetAuthorId(),
		UpdatedAt:      resp.GetUpdatedAt().AsTime(),
		CreatedAt:      resp.GetCreatedAt().AsTime(),
	}, nil
}

func (c *ArticleClient) GetArticlesByAuthorID(ctx context.Context, authorID int64) ([]model.Article, error) {
	resp, err := c.api.GetArticlesByAuthorID(ctx, &articleServicev1.GetArticlesByAuthorIDRequest{
		AuthorId: authorID,
	})

	if err != nil {
		_, ok := status.FromError(err)
		if ok {
			return nil, ErrInternalGRPCServer
		}
		return nil, err
	}

	articles := make([]model.Article, 0, len(resp.GetArticles()))
	for _, article := range resp.GetArticles() {
		articles = append(articles, model.Article{
			ID:         article.GetArticleId(),
			Title:      article.GetTitle(),
			AuthorName: article.GetAuthorName(),
			CreatedAt:  article.GetCreatedAt().AsTime(),
			UpdatedAt:  article.GetUpdatedAt().AsTime(),
		})
	}

	return articles, nil
}

func (c *ArticleClient) GetAllArticles(ctx context.Context) ([]model.Article, error) {
	resp, err := c.api.GetAllArticle(ctx, &articleServicev1.GetAllArticleRequest{})

	if err != nil {
		_, ok := status.FromError(err)
		if ok {
			return nil, ErrInternalGRPCServer
		}
		return nil, err
	}

	articles := make([]model.Article, 0, len(resp.GetArticles()))
	for _, article := range resp.GetArticles() {
		articles = append(articles, model.Article{
			ID:         article.GetArticleId(),
			Title:      article.GetTitle(),
			AuthorName: article.GetAuthorName(),
			CreatedAt:  article.GetCreatedAt().AsTime(),
			UpdatedAt:  article.GetUpdatedAt().AsTime(),
		})
	}

	return articles, nil
}

func (c *ArticleClient) UpdateArticle(ctx context.Context, input dto.UpdateArticleInput) error {
	_, err := c.api.UpdateArticle(ctx, &articleServicev1.UpdateArticleRequest{
		ArticleId: input.ID,
		Title:     input.Title,
		Content:   input.Content,
	})

	if err != nil {
		_, ok := status.FromError(err)
		if ok {
			return ErrInternalGRPCServer
		}
		return err
	}

	return nil
}

func (c *ArticleClient) DeleteArticle(ctx context.Context, id string) error {
	_, err := c.api.DeleteArticle(ctx, &articleServicev1.DeleteArticleRequest{
		ArticleId: id,
	})

	if err != nil {
		_, ok := status.FromError(err)
		if ok {
			return ErrInternalGRPCServer
		}
		return err
	}

	return nil
}
