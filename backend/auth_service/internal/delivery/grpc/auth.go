package grpc

import (
	"context"
	"errors"
	"github.com/Nikita-Mihailuk/goboard/backend/auth_service/internal/domain/dto"
	"github.com/Nikita-Mihailuk/goboard/backend/auth_service/internal/service/auth"
	authServicev1 "github.com/Nikita-Mihailuk/protos_goboard/gen/go/auth_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *serverGRPC) Register(ctx context.Context, req *authServicev1.RegisterRequest) (*authServicev1.RegisterResponse, error) {
	input, err := validateRegisterRequest(req)
	if err != nil {
		return nil, err
	}

	err = s.authService.Register(ctx, input)
	if err != nil {
		if errors.Is(err, auth.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		if errors.Is(err, auth.ErrInternalUserService) {
			return nil, status.Error(codes.Internal, "internal user service error")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &authServicev1.RegisterResponse{}, nil
}

func (s *serverGRPC) Login(ctx context.Context, req *authServicev1.LoginRequest) (*authServicev1.LoginResponse, error) {
	if req.GetEmail() == "" || req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "missing fields")
	}

	tokens, err := s.authService.Login(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid credentials")
		}
		if errors.Is(err, auth.ErrInternalUserService) {
			return nil, status.Error(codes.Internal, "internal user service error")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &authServicev1.LoginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (s *serverGRPC) Refresh(ctx context.Context, req *authServicev1.RefreshRequest) (*authServicev1.RefreshResponse, error) {
	if req.GetRefreshToken() == "" {
		return nil, status.Error(codes.InvalidArgument, "missing fields")
	}

	tokens, err := s.authService.RefreshTokens(ctx, req.GetRefreshToken())
	if err != nil {
		if errors.Is(err, auth.ErrInvalidRefreshToken) {
			return nil, status.Error(codes.InvalidArgument, "invalid refresh token")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &authServicev1.RefreshResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (s *serverGRPC) Logout(ctx context.Context, req *authServicev1.LogoutRequest) (*authServicev1.LogoutResponse, error) {
	if req.GetUserId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "missing fields")
	}
	err := s.authService.Logout(ctx, req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &authServicev1.LogoutResponse{}, nil
}

func validateRegisterRequest(req *authServicev1.RegisterRequest) (dto.RegisterUserInput, error) {
	if req.GetName() == "" || req.GetEmail() == "" || req.GetPassword() == "" {
		return dto.RegisterUserInput{}, status.Errorf(codes.InvalidArgument, "missing fields")
	}

	return dto.RegisterUserInput{
		Name:     req.GetName(),
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}, nil
}
