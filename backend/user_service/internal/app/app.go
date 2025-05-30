package app

import (
	grpcApp "github.com/Nikita-Mihailuk/goboard/backend/user_service/internal/app/grpc"
	"github.com/Nikita-Mihailuk/goboard/backend/user_service/internal/config"
	"github.com/Nikita-Mihailuk/goboard/backend/user_service/internal/service/user"
	"github.com/Nikita-Mihailuk/goboard/backend/user_service/internal/storage/postgres"
	"go.uber.org/zap"
)

type App struct {
	GRPCServer *grpcApp.App
}

func NewApp(log *zap.Logger, cfg *config.Config) *App {
	storage := postgres.NewStorage(cfg)

	service := user.NewUserService(log, storage, storage, storage)

	gRPCApp := grpcApp.NewApp(log, service, cfg.GRPCServer.Port)
	return &App{
		GRPCServer: gRPCApp,
	}
}
