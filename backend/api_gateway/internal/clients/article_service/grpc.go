package article_service

import (
	"context"
	articleServicev1 "github.com/Nikita-Mihailuk/protos_goboard/gen/go/article_service"
	"google.golang.org/grpc"
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
