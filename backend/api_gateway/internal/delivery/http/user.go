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
	userGroup.Get("/:id", h.getUserByID)
	userGroup.Patch("/:id", h.updateUser)

}

func (h *Handler) getUserByID(c fiber.Ctx) error {
	userId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid user id")
	}

	outputUser, err := h.userServiceClient.GetUserByID(c.Context(), int64(userId))
	if err != nil {
		if errors.Is(err, user_service.ErrUserNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "user not found")
		}
		if errors.Is(err, user_service.ErrInternalGRPC) {
			return fiber.NewError(fiber.StatusInternalServerError, "internal gRPC server error")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "internal error")
	}

	return c.JSON(outputUser)
}

func (h *Handler) updateUser(c fiber.Ctx) error {
	userId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid user id")
	}

	file, _ := c.FormFile("photo")
	inputUser := dto.UpdateUserInput{
		ID:              int64(userId),
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
