package app

import (
	grpcApp "github.com/Nikita-Mihailuk/goboard/backend/article_service/internal/app/grpc"
	"github.com/Nikita-Mihailuk/goboard/backend/article_service/internal/config"
	"github.com/Nikita-Mihailuk/goboard/backend/article_service/internal/infrastructure/kafka"
	"github.com/Nikita-Mihailuk/goboard/backend/article_service/internal/service/article"
	"github.com/Nikita-Mihailuk/goboard/backend/article_service/internal/storage/mongo"
	"go.uber.org/zap"
)

type App struct {
	GRPCServer    *grpcApp.App
	KafkaConsumer *kafka.Consumer
}

func NewApp(log *zap.Logger, cfg *config.Config) *App {
	storage := mongo.NewStorage(cfg)

	articleService := article.NewArticleService(log, storage, storage, storage, storage)

	gRPCApp := grpcApp.NewApp(log, articleService, cfg.GRPCServer.Port)

	messageHandler := kafka.NewHandler(storage)

	kafkaConsumer, err := kafka.NewConsumer(messageHandler, cfg.Kafka.Address, cfg.Kafka.UserServiceTopic, cfg.Kafka.ConsumerGroup, log)
	if err != nil {
		panic(err)
	}

	return &App{
		GRPCServer:    gRPCApp,
		KafkaConsumer: kafkaConsumer,
	}
}

func (a *App) Stop() {
	a.GRPCServer.Stop()
	if err := a.KafkaConsumer.Stop(); err != nil {
		panic(err)
	}
}
