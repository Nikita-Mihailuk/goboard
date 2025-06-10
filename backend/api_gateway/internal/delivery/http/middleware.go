package http

import (
	"github.com/gofiber/fiber/v3"
	"strings"
)

func (h *Handler) authMiddleware(c fiber.Ctx) error {
	accessToken := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
	if accessToken == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "access token not found")
	}

	claims, err := h.tokenManager.ParseAccessToken(accessToken)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid access token")
	}

	c.Locals("userID", claims.Subject)
	return c.Next()
}
