package auth

import (
	"context"
	"github.com/Nikita-Mihailuk/goboard/backend/auth_service/internal/clients/user_service"
	"github.com/Nikita-Mihailuk/goboard/backend/auth_service/pkg/auth"
	"go.uber.org/zap"
	"time"
)

type AuthService struct {
	log                  *zap.Logger
	tokenManager         *auth.TokenManager
	userServiceClient    *user_service.UserClient
	accessTokenTTL       time.Duration
	refreshTokenSaver    RefreshTokenSaver
	refreshTokenProvider RefreshTokenProvider
	refreshTokenDeleter  RefreshTokenDeleter
}

func NewArticleService(
	log *zap.Logger,
	tokenManager *auth.TokenManager,
	userServiceClient *user_service.UserClient,
	accessTokenTTL time.Duration,
	refreshTokenSaver RefreshTokenSaver,
	refreshTokenProvider RefreshTokenProvider,
	refreshTokenDeleter RefreshTokenDeleter,
) *AuthService {

	return &AuthService{
		log:                  log,
		tokenManager:         tokenManager,
		userServiceClient:    userServiceClient,
		accessTokenTTL:       accessTokenTTL,
		refreshTokenSaver:    refreshTokenSaver,
		refreshTokenProvider: refreshTokenProvider,
		refreshTokenDeleter:  refreshTokenDeleter,
	}
}

type RefreshTokenSaver interface {
	SetRefreshToken(ctx context.Context, userID int64, refreshToken string) error
}

type RefreshTokenProvider interface {
	GetRefreshToken(ctx context.Context, userID int64) (string, error)
}

type RefreshTokenDeleter interface {
	DeleteRefreshToken(ctx context.Context, userID int64) error
}
