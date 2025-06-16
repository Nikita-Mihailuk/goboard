package http

import (
	"errors"
	"github.com/Nikita-Mihailuk/goboard/backend/api_gateway/internal/clients/comment_service"
	"github.com/Nikita-Mihailuk/goboard/backend/api_gateway/internal/domain/dto"
	"github.com/gofiber/fiber/v3"
)

func (h *Handler) RegisterCommentRouts(router fiber.Router) {
	commentGroup := router.Group("/comments", h.authMiddleware)

	commentGroup.Post("", h.createComment)
	commentGroup.Get("/article/:articleID", h.getCommentsByArticleID)
	commentGroup.Patch("/:id", h.updateComment)
	commentGroup.Delete("/:id", h.deleteComment)
}

func (h *Handler) createComment(c fiber.Ctx) error {
	var input dto.CreateCommentInput
	if err := c.Bind().JSON(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request")
	}

	if err := h.commentServiceClient.CreateComment(c.Context(), input); err != nil {
		if errors.Is(err, comment_service.ErrInternalHTTP) {
			return fiber.NewError(fiber.StatusInternalServerError, "internal http server error")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "internal error")
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (h *Handler) getCommentsByArticleID(c fiber.Ctx) error {
	articleID := c.Params("articleID")

	comments, err := h.commentServiceClient.GetCommentsByArticleID(c.Context(), articleID)
	if err != nil {
		if errors.Is(err, comment_service.ErrInternalHTTP) {
			return fiber.NewError(fiber.StatusInternalServerError, "internal http server error")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}

	return c.JSON(comments)
}

func (h *Handler) updateComment(c fiber.Ctx) error {
	commentID := c.Params("id")

	var input dto.UpdateCommentInput
	if err := c.Bind().JSON(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request")
	}

	input.ID = commentID
	if err := h.commentServiceClient.UpdateComment(c.Context(), input); err != nil {
		if errors.Is(err, comment_service.ErrInternalHTTP) {
			return fiber.NewError(fiber.StatusInternalServerError, "internal http server error")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "internal error")
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) deleteComment(c fiber.Ctx) error {
	commentID := c.Params("id")

	if err := h.commentServiceClient.DeleteComment(c.Context(), commentID); err != nil {
		if errors.Is(err, comment_service.ErrInternalHTTP) {
			return fiber.NewError(fiber.StatusInternalServerError, "internal http server error")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "internal error")
	}

	return c.SendStatus(fiber.StatusNoContent)
}
