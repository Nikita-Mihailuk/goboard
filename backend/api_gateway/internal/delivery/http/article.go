package http

import (
	"errors"
	"github.com/Nikita-Mihailuk/goboard/backend/api_gateway/internal/clients/article_service"
	"github.com/Nikita-Mihailuk/goboard/backend/api_gateway/internal/domain/dto"
	"github.com/gofiber/fiber/v3"
	"strconv"
)

func (h *Handler) RegisterArticleRouts(router fiber.Router) {
	articleGroup := router.Group("/articles", h.authMiddleware)

	articleGroup.Post("", h.createArticle)
	articleGroup.Get("/:id", h.getArticleByID)
	articleGroup.Get("/author", h.getArticlesByAuthorID)
	articleGroup.Get("", h.getAllArticles)
	articleGroup.Patch("/:id", h.updateArticle)
	articleGroup.Delete("/:id", h.deleteArticle)
}

func (h *Handler) createArticle(c fiber.Ctx) error {
	var article dto.CreateArticleInput
	if err := c.Bind().JSON(&article); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request")
	}

	err := h.articleServiceClient.CreateArticle(c.Context(), article)
	if err != nil {
		if errors.Is(err, article_service.ErrInternalGRPCServer) {
			return fiber.NewError(fiber.StatusInternalServerError, "internal gRPC server error")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "internal error")
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (h *Handler) getArticleByID(c fiber.Ctx) error {
	articleID := c.Params("id")
	article, err := h.articleServiceClient.GetArticleByID(c.Context(), articleID)
	if err != nil {
		if errors.Is(err, article_service.ErrArticleNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "article not found")
		}
		if errors.Is(err, article_service.ErrInternalGRPCServer) {
			return fiber.NewError(fiber.StatusInternalServerError, "internal gRPC server error")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "internal error")
	}

	return c.JSON(article)
}

func (h *Handler) getArticlesByAuthorID(c fiber.Ctx) error {
	userIDstr, ok := c.Locals("userID").(string)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "user id not found")
	}

	userID, err := strconv.ParseInt(userIDstr, 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "invalid user id")
	}

	articles, err := h.articleServiceClient.GetArticlesByAuthorID(c.Context(), userID)

	if err != nil {
		if errors.Is(err, article_service.ErrInternalGRPCServer) {
			return fiber.NewError(fiber.StatusInternalServerError, "internal gRPC server error")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "internal error")
	}

	return c.JSON(articles)
}

func (h *Handler) getAllArticles(c fiber.Ctx) error {
	articles, err := h.articleServiceClient.GetAllArticles(c.Context())

	if err != nil {
		if errors.Is(err, article_service.ErrInternalGRPCServer) {
			return fiber.NewError(fiber.StatusInternalServerError, "internal gRPC server error")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "internal error")
	}

	return c.JSON(articles)
}

func (h *Handler) updateArticle(c fiber.Ctx) error {
	articleID := c.Params("id")

	var input dto.UpdateArticleInput
	if err := c.Bind().JSON(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request")
	}
	input.ID = articleID

	err := h.articleServiceClient.UpdateArticle(c.Context(), input)
	if err != nil {
		if errors.Is(err, article_service.ErrInternalGRPCServer) {
			return fiber.NewError(fiber.StatusInternalServerError, "internal gRPC server error")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "internal error")
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) deleteArticle(c fiber.Ctx) error {
	articleID := c.Params("id")
	err := h.articleServiceClient.DeleteArticle(c.Context(), articleID)
	if err != nil {
		if errors.Is(err, article_service.ErrInternalGRPCServer) {
			return fiber.NewError(fiber.StatusInternalServerError, "internal gRPC server error")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "internal error")
	}

	return c.SendStatus(fiber.StatusNoContent)
}
