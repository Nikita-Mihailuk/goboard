package comment_service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Nikita-Mihailuk/goboard/backend/api_gateway/internal/domain/dto"
	"github.com/Nikita-Mihailuk/goboard/backend/api_gateway/internal/domain/model"
	"net/http"
)

type CommentClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewCommentClient(baseURL string) *CommentClient {
	return &CommentClient{
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

func (c *CommentClient) CreateComment(ctx context.Context, input dto.CreateCommentInput) error {
	url := fmt.Sprintf("%s/comments", c.baseURL)
	body, _ := json.Marshal(input)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return ErrInternalHTTP
	}
	return nil
}

func (c *CommentClient) GetCommentsByArticleID(ctx context.Context, articleID string) ([]model.Comment, error) {
	url := fmt.Sprintf("%s/comments/article/%s", c.baseURL, articleID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, ErrInternalHTTP
	}

	var comments []model.Comment
	if err = json.NewDecoder(resp.Body).Decode(&comments); err != nil {
		return nil, err
	}

	return comments, nil
}

func (c *CommentClient) UpdateComment(ctx context.Context, input dto.UpdateCommentInput) error {
	url := fmt.Sprintf("%s/comments/%s", c.baseURL, input.ID)
	body, _ := json.Marshal(input)

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return ErrInternalHTTP
	}

	return nil
}

func (c *CommentClient) DeleteComment(ctx context.Context, id string) error {
	url := fmt.Sprintf("%s/comments/%s", c.baseURL, id)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return ErrInternalHTTP
	}

	return nil
}
