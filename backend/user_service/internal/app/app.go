package app

import (
	"github.com/Nikita-Mihailuk/goboard/backend/user_service/internal/app/http"
	"github.com/Nikita-Mihailuk/goboard/backend/user_service/internal/config"

	"go.uber.org/zap"
)

type App struct {
	HTTPServer *http.App
}

func NewApp(log *zap.Logger, cfg *config.Config) *App {

	// TODO
	return &App{
		HTTPServer: ,
	}
}
