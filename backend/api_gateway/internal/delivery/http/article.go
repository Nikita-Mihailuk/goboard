package http

import "github.com/gofiber/fiber/v3"

func (h *Handler) RegisterArticleRouts(router fiber.Router) {
	articleGroup := router.Group("/articles")

}
