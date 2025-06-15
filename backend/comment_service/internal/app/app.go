package app

import (
	httpApp "github.com/Nikita-Mihailuk/goboard/backend/comment_service/internal/app/http"
	"github.com/Nikita-Mihailuk/goboard/backend/comment_service/internal/config"
	"github.com/Nikita-Mihailuk/goboard/backend/comment_service/internal/delivery/http"
	"github.com/Nikita-Mihailuk/goboard/backend/comment_service/internal/infrastructure/kafka"
	"github.com/Nikita-Mihailuk/goboard/backend/comment_service/internal/service/comment"
	"github.com/Nikita-Mihailuk/goboard/backend/comment_service/internal/storage/mongo"
	"go.uber.org/zap"
)

type App struct {
	HTTPServer    *httpApp.App
	KafkaConsumer *kafka.Consumer
}

func NewApp(log *zap.Logger, cfg *config.Config) *App {
	storage := mongo.NewStorage(cfg)

	commentService := comment.NewCommentService(log, storage, storage, storage, storage)

	handler := http.NewHandler(commentService)

	httpHandler := handler.InitRoutes()

	httpServer := httpApp.NewApp(log, cfg.HTTPServer.Port, httpHandler)

	messageHandler := kafka.NewHandler(storage)

	kafkaConsumer, err := kafka.NewConsumer(messageHandler, cfg.Kafka.Address, cfg.Kafka.UserServiceTopic, cfg.Kafka.ConsumerGroup, log)
	if err != nil {
		panic(err)
	}

	return &App{
		HTTPServer:    httpServer,
		KafkaConsumer: kafkaConsumer,
	}
}

func (a *App) Stop() {
	a.HTTPServer.Stop()
	if err := a.KafkaConsumer.Stop(); err != nil {
		panic(err)
	}
}
