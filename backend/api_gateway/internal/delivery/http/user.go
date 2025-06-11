package http

import (
	"errors"
	"github.com/Nikita-Mihailuk/goboard/backend/api_gateway/internal/clients/user_service"
	"github.com/Nikita-Mihailuk/goboard/backend/api_gateway/internal/domain/dto"
	"github.com/gofiber/fiber/v3"
	"strconv"
)

func (h *Handler) RegisterUserRouts(router fiber.Router) {

	userGroup := router.Group("/users", h.authMiddleware)
	userGroup.Get("/", h.getUserByID)
	userGroup.Patch("/", h.updateUser)

}

func (h *Handler) getUserByID(c fiber.Ctx) error {
	userIDstr, ok := c.Locals("userID").(string)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "user id not found")
	}

	userID, err := strconv.ParseInt(userIDstr, 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "invalid user id")
	}

	outputUser, err := h.userServiceClient.GetUserByID(c.Context(), userID)
	if err != nil {
		if errors.Is(err, user_service.ErrUserNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "user not found")
		}
		if errors.Is(err, user_service.ErrInternalGRPC) {
			return fiber.NewError(fiber.StatusInternalServerError, "internal gRPC server error")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "internal error")
	}

	return c.JSON(fiber.Map{
		"user_id": userID,
		"user":    outputUser,
	})
}

func (h *Handler) updateUser(c fiber.Ctx) error {
	userIDstr, ok := c.Locals("userID").(string)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "user id not found")
	}

	userID, err := strconv.ParseInt(userIDstr, 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "invalid user id")
	}

	file, _ := c.FormFile("photo")
	inputUser := dto.UpdateUserInput{
		ID:              userID,
		CurrentPassword: c.FormValue("current_password"),
		NewPassword:     c.FormValue("new_password"),
		Name:            c.FormValue("new_name"),
		FileHeader:      file,
	}

	err = h.userServiceClient.UpdateUser(c.Context(), inputUser)
	if err != nil {
		if errors.Is(err, user_service.ErrUserNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "user not found")
		}
		if errors.Is(err, user_service.ErrInvalidPassword) {
			return fiber.NewError(fiber.StatusBadRequest, "invalid password")
		}
		if errors.Is(err, user_service.ErrInternalGRPC) {
			return fiber.NewError(fiber.StatusInternalServerError, "internal gRPC server error")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "internal error")
	}

	return c.SendStatus(fiber.StatusNoContent)
}
