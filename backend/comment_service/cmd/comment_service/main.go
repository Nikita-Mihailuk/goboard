package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Nikita-Mihailuk/goboard/backend/comment_service/internal/app"
	"github.com/Nikita-Mihailuk/goboard/backend/comment_service/internal/config"
	"github.com/Nikita-Mihailuk/goboard/backend/comment_service/pkg/logging"
)

func main() {
	cfg := config.GetConfig()
	log := logging.GetLogger(cfg.Env)

	application := app.NewApp(log, cfg)
	go application.HTTPServer.MustRun()
	//go application.KafkaConsumer.Start()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	application.Stop()
	log.Info("application stopped")
}
