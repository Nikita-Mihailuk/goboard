package app

import (
	grpcApp "github.com/Nikita-Mihailuk/goboard/backend/auth_service/internal/app/grpc"
	"github.com/Nikita-Mihailuk/goboard/backend/auth_service/internal/config"
	"github.com/Nikita-Mihailuk/goboard/backend/auth_service/internal/service/auth"
	"github.com/Nikita-Mihailuk/goboard/backend/auth_service/internal/storage/redis"
	authManager "github.com/Nikita-Mihailuk/goboard/backend/auth_service/pkg/auth"
	"go.uber.org/zap"
)

type App struct {
	GRPCServer *grpcApp.App
}

func NewApp(log *zap.Logger, cfg *config.Config) *App {
	storage := redis.NewStorage(cfg, cfg.Auth.RefreshTokenTTL)

	tokenManager, err := authManager.NewManager(cfg.Auth.SecretKey)
	if err != nil {
		panic(err)
	}
	articleService := auth.NewArticleService(log, tokenManager, cfg.Auth.AccessTokenTTL, storage, storage, storage)

	gRPCApp := grpcApp.NewApp(log, articleService, cfg.GRPCServer.Port)
	return &App{
		GRPCServer: gRPCApp,
	}
}
