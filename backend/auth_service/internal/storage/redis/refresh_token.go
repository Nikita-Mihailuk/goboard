package redis

import (
	"context"
	"fmt"
)

func (s *Storage) SetRefreshToken(ctx context.Context, userID int64, refreshToken string) error {
	key := fmt.Sprintf("user:%d", userID)
	return s.db.Set(ctx, key, refreshToken, s.refreshTokenTTL).Err()
}

func (s *Storage) GetRefreshToken(ctx context.Context, userID int64) (string, error) {
	key := fmt.Sprintf("user:%d", userID)

	token, err := s.db.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Storage) DeleteRefreshToken(ctx context.Context, userID int64) error {
	key := fmt.Sprintf("user:%d", userID)
	return s.db.Del(ctx, key).Err()
}
