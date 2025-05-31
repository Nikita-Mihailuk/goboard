package http

import "github.com/gofiber/fiber/v3"

func (h *Handler) RegisterUserRouts(router fiber.Router) {
	userGroup := router.Group("/users")

}
