package http

import (
	"github.com/Nikita-Mihailuk/goboard/backend/api_gateway/internal/clients/article_service"
	"github.com/Nikita-Mihailuk/goboard/backend/api_gateway/internal/clients/auth_service"
	"github.com/Nikita-Mihailuk/goboard/backend/api_gateway/internal/clients/comment_service"
	"github.com/Nikita-Mihailuk/goboard/backend/api_gateway/internal/clients/user_service"
	"github.com/Nikita-Mihailuk/goboard/backend/api_gateway/pkg/auth"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	userServiceClient    *user_service.UserClient
	articleServiceClient *article_service.ArticleClient
	authServiceClient    *auth_service.AuthClient
	commentServiceClient *comment_service.CommentClient
	tokenManager         *auth.Manager
}

func NewHandler(
	userServiceClient *user_service.UserClient,
	articleServiceClient *article_service.ArticleClient,
	authServiceClient *auth_service.AuthClient,
	commentServiceClient *comment_service.CommentClient,
	tokenManager *auth.Manager,
) *Handler {
	return &Handler{
		userServiceClient:    userServiceClient,
		articleServiceClient: articleServiceClient,
		authServiceClient:    authServiceClient,
		commentServiceClient: commentServiceClient,
		tokenManager:         tokenManager,
	}
}

func (h *Handler) InitRoutes(router fiber.Router) {
	router.Get("/ping", func(ctx fiber.Ctx) error {
		return ctx.SendString("pong")
	})

	h.RegisterUserRouts(router)
	h.RegisterArticleRouts(router)
	h.RegisterAuthRouts(router)
	h.RegisterCommentRouts(router)
}
