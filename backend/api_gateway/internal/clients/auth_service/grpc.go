package auth_service

import (
	"context"
	"github.com/Nikita-Mihailuk/goboard/backend/api_gateway/internal/domain/dto"
	authServicev1 "github.com/Nikita-Mihailuk/protos_goboard/gen/go/auth_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthClient struct {
	api authServicev1.AuthClient
}

func NewAuthClient(ctx context.Context, addr string) (*AuthClient, error) {
	cc, err := grpc.DialContext(ctx, addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &AuthClient{api: authServicev1.NewAuthClient(cc)}, nil
}

func (c *AuthClient) Register(ctx context.Context, input dto.RegisterUserInput) error {
	_, err := c.api.Register(ctx, &authServicev1.RegisterRequest{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	})

	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.AlreadyExists:
				return ErrUserExists
			default:
				return ErrInternalGRPC
			}
		}
		return err
	}

	return nil
}

func (c *AuthClient) Login(ctx context.Context, input dto.LoginUserInput) (dto.Tokens, error) {
	resp, err := c.api.Login(ctx, &authServicev1.LoginRequest{
		Email:    input.Email,
		Password: input.Password,
	})

	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.InvalidArgument:
				return dto.Tokens{}, ErrInvalidCredentials
			default:
				return dto.Tokens{}, ErrInternalGRPC
			}
		}
		return dto.Tokens{}, err
	}

	return dto.Tokens{
		AccessToken:  resp.GetAccessToken(),
		RefreshToken: resp.GetRefreshToken(),
	}, nil
}

func (c *AuthClient) Refresh(ctx context.Context, refreshToken string) (dto.Tokens, error) {
	resp, err := c.api.Refresh(ctx, &authServicev1.RefreshRequest{
		RefreshToken: refreshToken,
	})
	if err != nil {
		_, ok := status.FromError(err)
		if ok {
			return dto.Tokens{}, ErrInternalGRPC
		}
		return dto.Tokens{}, err
	}

	return dto.Tokens{
		AccessToken:  resp.GetAccessToken(),
		RefreshToken: resp.GetRefreshToken(),
	}, nil
}

func (c *AuthClient) Logout(ctx context.Context, userID int64) error {
	_, err := c.api.Logout(ctx, &authServicev1.LogoutRequest{
		UserId: userID,
	})
	if err != nil {
		_, ok := status.FromError(err)
		if ok {
			return ErrInternalGRPC
		}
		return err
	}

	return nil
}
