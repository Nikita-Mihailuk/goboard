package redis

import "context"

func (s *Storage) SetRefreshToken(ctx context.Context, refreshToken string) error {
	panic("implement me")
}

func (s *Storage) GetRefreshToken(ctx context.Context, userID int64) (string, error) {
	panic("implement me")
}

func (s *Storage) DeleteRefreshToken(ctx context.Context, userID int64) error {
	panic("implement me")
}
