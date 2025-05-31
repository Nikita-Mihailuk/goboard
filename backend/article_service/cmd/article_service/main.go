package main

import (
	"github.com/Nikita-Mihailuk/goboard/backend/article_service/internal/app"
	"github.com/Nikita-Mihailuk/goboard/backend/article_service/internal/config"
	"github.com/Nikita-Mihailuk/goboard/backend/article_service/pkg/logging"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.GetConfig()
	log := logging.GetLogger(cfg.Env)

	application := app.NewApp(log, cfg)
	go application.GRPCServer.MustRun()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	application.GRPCServer.Stop()
	log.Info("application stopped")
}
