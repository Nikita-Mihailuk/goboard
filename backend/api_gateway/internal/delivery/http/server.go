package http

import (
	"github.com/Nikita-Mihailuk/goboard/backend/api_gateway/internal/clients/article_service"
	"github.com/Nikita-Mihailuk/goboard/backend/api_gateway/internal/clients/user_service"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	userServiceClient    *user_service.UserClient
	articleServiceClient *article_service.ArticleClient
}

func NewHandler(userServiceClient *user_service.UserClient, articleServiceClient *article_service.ArticleClient) *Handler {
	return &Handler{
		userServiceClient:    userServiceClient,
		articleServiceClient: articleServiceClient,
	}
}

func (h *Handler) InitRoutes(router fiber.Router) {
	router.Get("/ping", func(ctx fiber.Ctx) error {
		return ctx.SendString("pong")
	})

	h.RegisterUserRouts(router)
	h.RegisterArticleRouts(router)
}
