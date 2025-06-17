package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/Nikita-Mihailuk/goboard/backend/api_gateway/internal/domain/model"
	"strconv"
	"time"
)

var (
	ErrArticleNotFound = errors.New("article not found")
)

func (c *Cache) SetArticle(ctx context.Context, article model.Article) error {
	articleMap := structToMap(&article)

	key := fmt.Sprintf("article:%s", article.ID)

	if err := c.redis.HMSet(ctx, key, articleMap).Err(); err != nil {
		return err
	}

	return c.redis.Expire(ctx, key, c.cacheTTL).Err()
}

func (c *Cache) GetArticle(ctx context.Context, id string) (model.Article, error) {
	key := fmt.Sprintf("article:%s", id)

	articleMap, err := c.redis.HGetAll(ctx, key).Result()
	if err != nil {
		return model.Article{}, err
	}

	if len(articleMap) == 0 {
		return model.Article{}, ErrArticleNotFound
	}

	authorID, _ := strconv.ParseInt(articleMap["author_id"], 10, 64)
	createdAt, _ := time.Parse(time.RFC3339, articleMap["created_at"])
	updatedAt, _ := time.Parse(time.RFC3339, articleMap["updated_at"])

	return model.Article{
		ID:             articleMap["id"],
		Title:          articleMap["title"],
		Content:        articleMap["content"],
		AuthorName:     articleMap["author_name"],
		AuthorID:       authorID,
		AuthorPhotoURL: articleMap["author_photo_url"],
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}, nil
}
