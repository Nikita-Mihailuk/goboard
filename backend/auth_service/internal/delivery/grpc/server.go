package grpc

import (
	"context"
	"github.com/Nikita-Mihailuk/goboard/backend/auth_service/internal/domain/dto"
	authServicev1 "github.com/Nikita-Mihailuk/protos_goboard/gen/go/auth_service"
	"google.golang.org/grpc"
)

type AuthService interface {
	Register(ctx context.Context, input dto.RegisterUserInput) error
	Login(ctx context.Context, email, password string) (dto.Tokens, error)
	RefreshTokens(ctx context.Context, refreshToken string) (dto.Tokens, error)
	Logout(ctx context.Context, userID int64) error
}

type serverGRPC struct {
	authService AuthService
	authServicev1.UnimplementedAuthServer
}

func RegisterGRPCServer(grpcServer *grpc.Server, authService AuthService) {
	authServicev1.RegisterAuthServer(grpcServer, &serverGRPC{authService: authService})
}
