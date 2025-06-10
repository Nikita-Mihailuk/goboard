package http

import (
	"errors"
	"github.com/Nikita-Mihailuk/goboard/backend/api_gateway/internal/clients/auth_service"
	"github.com/Nikita-Mihailuk/goboard/backend/api_gateway/internal/domain/dto"
	"github.com/gofiber/fiber/v3"
	"strconv"
)

func (h *Handler) RegisterAuthRouts(router fiber.Router) {
	authGroup := router.Group("/auth")
	authGroup.Post("/login", h.loginUser)
	authGroup.Post("/register", h.registerUser)
	authGroup.Post("/refresh", h.refreshUser)
	authGroup.Delete("/logout", h.logoutUser, h.authMiddleware)
}

func (h *Handler) loginUser(c fiber.Ctx) error {
	var inputLogin dto.LoginUserInput
	if err := c.Bind().JSON(&inputLogin); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request")
	}

	tokens, err := h.authServiceClient.Login(c.Context(), inputLogin)
	if err != nil {
		if errors.Is(err, auth_service.ErrInvalidCredentials) {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
		}
		if errors.Is(err, auth_service.ErrInternalGRPC) {
			return fiber.NewError(fiber.StatusInternalServerError, "internal gRPC server error")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "internal error")
	}

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		HTTPOnly: true,
		Secure:   false, // true для https
	})

	return c.JSON(fiber.Map{
		"access_token": tokens.AccessToken,
	})
}

func (h *Handler) registerUser(c fiber.Ctx) error {
	var inputRegister dto.RegisterUserInput
	if err := c.Bind().JSON(&inputRegister); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request")
	}

	err := h.authServiceClient.Register(c.Context(), inputRegister)
	if err != nil {
		if errors.Is(err, auth_service.ErrUserExists) {
			return fiber.NewError(fiber.StatusConflict, "user already exists")
		}
		if errors.Is(err, auth_service.ErrInternalGRPC) {
			return fiber.NewError(fiber.StatusInternalServerError, "internal gRPC server error")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "internal error")
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (h *Handler) refreshUser(c fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return fiber.NewError(fiber.StatusBadRequest, "invalid refresh token")
	}

	newTokens, err := h.authServiceClient.Refresh(c.Context(), refreshToken)
	if err != nil {
		if errors.Is(err, auth_service.ErrInternalGRPC) {
			return fiber.NewError(fiber.StatusInternalServerError, "internal gRPC server error")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "internal error")
	}

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    newTokens.RefreshToken,
		HTTPOnly: true,
		Secure:   false, // true для https
	})

	return c.JSON(fiber.Map{
		"access_token": newTokens.AccessToken,
	})
}

func (h *Handler) logoutUser(c fiber.Ctx) error {
	userIDstr, ok := c.Locals("userID").(string)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "user id not found")
	}

	userID, err := strconv.ParseInt(userIDstr, 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "invalid user id")
	}
	err = h.authServiceClient.Logout(c.Context(), userID)
	if err != nil {
		if errors.Is(err, auth_service.ErrInternalGRPC) {
			return fiber.NewError(fiber.StatusInternalServerError, "internal gRPC server error")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "internal error")
	}

	c.ClearCookie("refresh_token")

	return c.SendStatus(fiber.StatusNoContent)
}
