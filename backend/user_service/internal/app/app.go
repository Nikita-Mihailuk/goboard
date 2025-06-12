package app

import (
	grpcApp "github.com/Nikita-Mihailuk/goboard/backend/user_service/internal/app/grpc"
	"github.com/Nikita-Mihailuk/goboard/backend/user_service/internal/config"
	"github.com/Nikita-Mihailuk/goboard/backend/user_service/internal/infrastructure/kafka"
	"github.com/Nikita-Mihailuk/goboard/backend/user_service/internal/service/user"
	"github.com/Nikita-Mihailuk/goboard/backend/user_service/internal/storage/postgres"
	"go.uber.org/zap"
)

type App struct {
	GRPCServer    *grpcApp.App
	kafkaProducer *kafka.Producer
}

func NewApp(log *zap.Logger, cfg *config.Config) *App {
	storage := postgres.NewStorage(cfg)

	kafkaProducer, err := kafka.NewProducer(cfg.Kafka.Address)
	if err != nil {
		panic(err)
	}

	service := user.NewUserService(log, storage, storage, storage, kafkaProducer, cfg.Kafka.ProducerTopic)

	gRPCApp := grpcApp.NewApp(log, service, cfg.GRPCServer.Port)
	return &App{
		GRPCServer:    gRPCApp,
		kafkaProducer: kafkaProducer,
	}
}

func (a *App) Stop() {
	a.kafkaProducer.Close()
	a.GRPCServer.Stop()
}
