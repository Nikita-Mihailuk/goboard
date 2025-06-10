package grpc

import (
	"context"
	authServicev1 "github.com/Nikita-Mihailuk/protos_goboard/gen/go/auth_service"
)

func (s *serverGRPC) Register(ctx context.Context, req *authServicev1.RegisterRequest) (*authServicev1.RegisterResponse, error) {
	panic("implement me")
}

func (s *serverGRPC) Login(ctx context.Context, req *authServicev1.LoginRequest) (*authServicev1.LoginResponse, error) {
	panic("implement me")
}

func (s *serverGRPC) Refresh(ctx context.Context, req *authServicev1.RefreshRequest) (*authServicev1.RefreshResponse, error) {
	panic("implement me")
}

func (s *serverGRPC) Logout(ctx context.Context, req *authServicev1.LogoutRequest) (*authServicev1.LogoutResponse, error) {
	panic("implement me")
}
