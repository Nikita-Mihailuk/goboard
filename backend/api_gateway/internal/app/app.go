package app

import (
	"context"
	"fmt"
	"github.com/Nikita-Mihailuk/goboard/backend/api_gateway/internal/app/http"
	"github.com/Nikita-Mihailuk/goboard/backend/api_gateway/internal/clients/article_service"
	"github.com/Nikita-Mihailuk/goboard/backend/api_gateway/internal/clients/auth_service"
	"github.com/Nikita-Mihailuk/goboard/backend/api_gateway/internal/clients/user_service"
	"github.com/Nikita-Mihailuk/goboard/backend/api_gateway/internal/config"
	"github.com/Nikita-Mihailuk/goboard/backend/api_gateway/pkg/auth"
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

	authServiceClient, err := auth_service.NewAuthClient(ctx, fmt.Sprintf("%s:%s", cfg.AuthService.Host, cfg.AuthService.Port))
	if err != nil {
		panic(err)
	}

	tokenManager, err := auth.NewManager(cfg.Auth.SecretKey)
	if err != nil {
		panic(err)
	}

	httpApp := http.NewApp(cfg.HTTPServer.Port, userServiceClient, articleServiceClient, authServiceClient, tokenManager)
	return &App{
		HTTPServer: httpApp,
	}
}
