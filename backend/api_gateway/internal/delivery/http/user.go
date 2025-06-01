package http

import (
	"errors"
	"github.com/Nikita-Mihailuk/goboard/backend/api_gateway/internal/clients/user_service"
	"github.com/Nikita-Mihailuk/goboard/backend/api_gateway/internal/domain/dto"
	"github.com/gofiber/fiber/v3"
	"strconv"
)

func (h *Handler) RegisterUserRouts(router fiber.Router) {
	// TODO take out in auth service
	authGroup := router.Group("/auth")
	authGroup.Post("/login", h.loginUser)
	authGroup.Post("/register", h.registerUser)

	userGroup := router.Group("/users")
	userGroup.Get("/:id", h.getUserByID)
	userGroup.Patch("/:id", h.updateUser)

}

// TODO take out in auth service
func (h *Handler) loginUser(c fiber.Ctx) error {
	var inputLogin dto.LoginUserInput
	if err := c.Bind().JSON(&inputLogin); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request")
	}

	userID, err := h.userServiceClient.LoginUser(c.Context(), inputLogin.Email, inputLogin.Password)
	if err != nil {
		if errors.Is(err, user_service.ErrInvalidCredentials) {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
		}
		if errors.Is(err, user_service.ErrInternalGRPC) {
			return fiber.NewError(fiber.StatusInternalServerError, "internal gRPC server error")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "internal error")
	}

	return c.JSON(fiber.Map{
		"user_id": userID,
	})
}

// TODO take out in auth service
func (h *Handler) registerUser(c fiber.Ctx) error {
	var inputRegister dto.CreateUserInput
	if err := c.Bind().JSON(&inputRegister); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request")
	}

	err := h.userServiceClient.CreateUser(c.Context(), inputRegister)
	if err != nil {
		if errors.Is(err, user_service.ErrUserExists) {
			return fiber.NewError(fiber.StatusConflict, "user already exists")
		}
		if errors.Is(err, user_service.ErrInternalGRPC) {
			return fiber.NewError(fiber.StatusInternalServerError, "internal gRPC server error")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "internal error")
	}

	return c.SendStatus(fiber.StatusCreated)
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
