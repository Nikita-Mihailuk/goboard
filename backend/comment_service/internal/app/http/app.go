package http

import (
	"context"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type App struct {
	log    *zap.Logger
	server *http.Server
}

func NewApp(log *zap.Logger, port string, handler http.Handler) *App {
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	return &App{
		log:    log,
		server: server,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	a.log.Info("HTTP server started on " + a.server.Addr)
	return a.server.ListenAndServe()
}

func (a *App) Stop() {
	a.log.Info("HTTP server stopped")
	if err := a.server.Shutdown(context.Background()); err != nil {
		a.log.Error("Failed to stop HTTP server", zap.Error(err))
	}
}
