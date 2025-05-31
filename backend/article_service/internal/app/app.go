package app

import (
	grpcApp "github.com/Nikita-Mihailuk/goboard/backend/article_service/internal/app/grpc"
	"github.com/Nikita-Mihailuk/goboard/backend/article_service/internal/config"
	"github.com/Nikita-Mihailuk/goboard/backend/article_service/internal/service/article"
	"github.com/Nikita-Mihailuk/goboard/backend/article_service/internal/storage/mongo"
	"go.uber.org/zap"
)

type App struct {
	GRPCServer *grpcApp.App
}

func NewApp(log *zap.Logger, cfg *config.Config) *App {
	storage := mongo.NewStorage(cfg)

	articleService := article.NewArticleService(log, storage, storage, storage, storage)

	gRPCApp := grpcApp.NewApp(log, articleService, cfg.GRPCServer.Port)
	return &App{
		GRPCServer: gRPCApp,
	}
}
