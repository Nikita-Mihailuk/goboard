package app

import (
	"context"
	"fmt"
	"github.com/Nikita-Mihailuk/goboard/backend/api_gateway/internal/app/http"
	"github.com/Nikita-Mihailuk/goboard/backend/api_gateway/internal/clients/article_service"
	"github.com/Nikita-Mihailuk/goboard/backend/api_gateway/internal/clients/user_service"
	"github.com/Nikita-Mihailuk/goboard/backend/api_gateway/internal/config"
)

type App struct {
	HTTPServer *http.App
}

func NewApp(cfg *config.Config) *App {

	ctx, cancel := context.WithTimeout(context.Background(), cfg.HTTPServer.Timeout)
	defer cancel()

	userServiceClient, err := user_service.NewUserClient(ctx, fmt.Sprintf("%s:%s", cfg.UserService.Host, cfg.UserService.Port))
	if err != nil {
		panic(err)
	}

	articleServiceClient, err := article_service.NewArticleClient(ctx, fmt.Sprintf("%s:%s", cfg.ArticleService.Host, cfg.ArticleService.Port))
	if err != nil {
		panic(err)
	}

	httpApp := http.NewApp(cfg.HTTPServer.Port, userServiceClient, articleServiceClient)
	return &App{
		HTTPServer: httpApp,
	}
}
