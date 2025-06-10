package user_service

import (
	"context"
	"github.com/Nikita-Mihailuk/goboard/backend/auth_service/internal/domain/dto"
	userServicev1 "github.com/Nikita-Mihailuk/protos_goboard/gen/go/user_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserClient struct {
	api userServicev1.UserClient
}

func NewUserClient(ctx context.Context, addr string) (*UserClient, error) {
	cc, err := grpc.DialContext(ctx, addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &UserClient{api: userServicev1.NewUserClient(cc)}, nil
}

func (c *UserClient) CreateUser(ctx context.Context, input dto.RegisterUserInput) error {
	_, err := c.api.CreateUser(ctx, &userServicev1.CreateUserRequest{
		Email:    input.Email,
		Name:     input.Name,
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

func (c *UserClient) LoginUser(ctx context.Context, email string) (dto.LoginUserOutput, error) {
	resp, err := c.api.LoginUser(ctx, &userServicev1.LoginUserRequest{
		Email: email,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				return dto.LoginUserOutput{}, ErrInvalidCredentials
			default:
				return dto.LoginUserOutput{}, ErrInternalGRPC
			}
		}
		return dto.LoginUserOutput{}, err
	}

	return dto.LoginUserOutput{
		ID:           resp.GetUserId(),
		PasswordHash: resp.GetPasswordHash(),
		Role:         resp.GetRole(),
	}, nil
}
