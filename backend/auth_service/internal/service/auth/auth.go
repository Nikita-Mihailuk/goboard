package auth

import (
	"context"
	"errors"
	"github.com/Nikita-Mihailuk/goboard/backend/auth_service/internal/clients/user_service"
	"github.com/Nikita-Mihailuk/goboard/backend/auth_service/internal/domain/dto"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"strconv"
)

func (s *AuthService) Register(ctx context.Context, input dto.RegisterUserInput) error {
	passHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		s.log.Error("failed to hash password", zap.Error(err))
		return err
	}
	input.Password = string(passHash)

	err = s.userServiceClient.CreateUser(ctx, input)
	if err != nil {
		if errors.Is(err, user_service.ErrUserExists) {
			s.log.Error("user already exists", zap.String("email", input.Email))
			return ErrUserExists
		}
		if errors.Is(err, user_service.ErrInternalGRPC) {
			s.log.Error("internal user service error", zap.Error(err))
			return ErrInternalUserService
		}
		s.log.Error("failed to create user", zap.Error(err))
		return err
	}

	s.log.Info("successfully created user", zap.String("email", input.Email))
	return nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (dto.Tokens, error) {
	outputLogin, err := s.userServiceClient.LoginUser(ctx, email)
	if err != nil {
		if errors.Is(err, user_service.ErrInvalidCredentials) {
			s.log.Error("invalid credentials", zap.String("email", email))
			return dto.Tokens{}, ErrInvalidCredentials
		}
		if errors.Is(err, user_service.ErrInternalGRPC) {
			s.log.Error("internal user service error", zap.Error(err))
			return dto.Tokens{}, ErrInternalUserService
		}
		s.log.Error("failed to login", zap.Error(err))
		return dto.Tokens{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(outputLogin.PasswordHash), []byte(password))
	if err != nil {
		s.log.Error("invalid credentials", zap.String("email", email))
		return dto.Tokens{}, ErrInvalidCredentials
	}

	refreshToken, err := s.tokenManager.NewRefreshToken(outputLogin.ID, outputLogin.Role)
	if err != nil {
		s.log.Error("failed to create refresh token", zap.Error(err))
		return dto.Tokens{}, err
	}

	err = s.refreshTokenSaver.SetRefreshToken(ctx, outputLogin.ID, refreshToken)
	if err != nil {
		s.log.Error("failed to save refresh token", zap.Error(err))
		return dto.Tokens{}, err
	}

	accessToken, err := s.tokenManager.NewJWT(strconv.FormatInt(outputLogin.ID, 10), outputLogin.Role, s.accessTokenTTL)
	if err != nil {
		s.log.Error("failed to create access token", zap.Error(err))
		return dto.Tokens{}, err
	}

	s.log.Info("successfully logged in", zap.Int64("userID", outputLogin.ID))
	return dto.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) RefreshTokens(ctx context.Context, refreshToken string) (dto.Tokens, error) {
	tokenData, err := s.tokenManager.ParseRefreshToken(refreshToken)
	if err != nil {
		s.log.Error("failed to parse refresh token", zap.Error(err))
		return dto.Tokens{}, err
	}

	redisToken, err := s.refreshTokenProvider.GetRefreshToken(ctx, tokenData.UserID)
	if err != nil {
		s.log.Error("failed to get refresh token", zap.Int64("user_id", tokenData.UserID))
		return dto.Tokens{}, err
	}

	if redisToken != refreshToken {
		s.log.Error("invalid refresh token", zap.Int64("user_id", tokenData.UserID))
		return dto.Tokens{}, ErrInvalidRefreshToken
	}

	newRefreshToken, err := s.tokenManager.NewRefreshToken(tokenData.UserID, tokenData.Role)
	if err != nil {
		s.log.Error("failed to create refresh token", zap.Error(err))
		return dto.Tokens{}, err
	}

	err = s.refreshTokenSaver.SetRefreshToken(ctx, tokenData.UserID, newRefreshToken)
	if err != nil {
		s.log.Error("failed to save refresh token", zap.Error(err))
		return dto.Tokens{}, err
	}

	newAccessToken, err := s.tokenManager.NewJWT(strconv.FormatInt(tokenData.UserID, 10), tokenData.Role, s.accessTokenTTL)
	if err != nil {
		s.log.Error("failed to create access token", zap.Error(err))
		return dto.Tokens{}, err
	}

	s.log.Info("successfully refresh tokens", zap.Int64("user_id", tokenData.UserID))
	return dto.Tokens{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, userID int64) error {
	err := s.refreshTokenDeleter.DeleteRefreshToken(ctx, userID)
	if err != nil {
		s.log.Error("failed to delete refresh token", zap.Int64("user_id", userID))
		return err
	}

	s.log.Info("successfully logout user", zap.Int64("user_id", userID))
	return nil
}
