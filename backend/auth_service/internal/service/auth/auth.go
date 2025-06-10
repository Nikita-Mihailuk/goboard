package auth

import (
	"context"
	"github.com/Nikita-Mihailuk/goboard/backend/auth_service/internal/domain/dto"
)

func (s *AuthService) Register(ctx context.Context, input dto.RegisterUserInput) error {
	panic("implement me")
}

func (s *AuthService) Login(ctx context.Context, email, password string) (dto.Tokens, error) {
	panic("implement me")
}

func (s *AuthService) RefreshTokens(ctx context.Context, refreshToken string) (dto.Tokens, error) {
	panic("implement me")
}

func (s *AuthService) Logout(ctx context.Context, userID int64) error {
	panic("implement me")
}
