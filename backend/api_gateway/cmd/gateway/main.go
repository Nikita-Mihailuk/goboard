package main

import (
	"github.com/Nikita-Mihailuk/goboard/backend/api_gateway/internal/app"
	"github.com/Nikita-Mihailuk/goboard/backend/api_gateway/internal/config"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.GetConfig()

	application := app.NewApp(cfg)

	go application.HTTPServer.MustRun()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	application.HTTPServer.Stop()
}
