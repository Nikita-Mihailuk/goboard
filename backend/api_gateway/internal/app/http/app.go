package http

import (
	"github.com/Nikita-Mihailuk/goboard/backend/api_gateway/internal/clients/comment_service"
	"time"

	"github.com/Nikita-Mihailuk/goboard/backend/api_gateway/internal/clients/article_service"
	"github.com/Nikita-Mihailuk/goboard/backend/api_gateway/internal/clients/auth_service"
	"github.com/Nikita-Mihailuk/goboard/backend/api_gateway/internal/clients/user_service"
	"github.com/Nikita-Mihailuk/goboard/backend/api_gateway/internal/delivery/http"
	"github.com/Nikita-Mihailuk/goboard/backend/api_gateway/pkg/auth"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/static"
)

type App struct {
	router *fiber.App
	port   string
}

func NewApp(
	port string,
	userServiceClient *user_service.UserClient,
	articleServiceClient *article_service.ArticleClient,
	authServiceClient *auth_service.AuthClient,
	commentServiceClient *comment_service.CommentClient,
	tokenManager *auth.Manager,
) *App {

	router := fiber.New(fiber.Config{
		ServerHeader: "Fiber",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	})
	router.Get("/static*", static.New("./static"))
	router.Use(logger.New(), cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost"},
		AllowCredentials: true,
		AllowHeaders:     []string{"Content-Type", "Authorization", "Cache-Control"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
	}))

	handler := http.NewHandler(userServiceClient, articleServiceClient, authServiceClient, commentServiceClient, tokenManager)
	handler.InitRoutes(router)

	return &App{
		router: router,
		port:   port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	if err := a.router.Listen(":" + a.port); err != nil {
		return err
	}
	return nil
}

func (a *App) Stop() {
	if err := a.router.Shutdown(); err != nil {
		panic(err)
	}
}
