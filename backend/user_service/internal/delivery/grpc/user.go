package grpc

import (
	"context"
	userServicev1 "github.com/Nikita-Mihailuk/protos_goboard/gen/go/user_service"
)

func (s *serverGRPC) CreateUser(context.Context, *userServicev1.CreateUserRequest) (*userServicev1.CreateUserResponse, error) {
	panic("implement me")
}

func (s *serverGRPC) LoginUser(context.Context, *userServicev1.LoginUserRequest) (*userServicev1.LoginUserResponse, error) {
	panic("implement me")
}

func (s *serverGRPC) GetUserByID(context.Context, *userServicev1.GetUserByIDRequest) (*userServicev1.GetUserByIDResponse, error) {
	panic("implement me")
}

func (s *serverGRPC) UpdateUser(context.Context, *userServicev1.UpdateUserRequest) (*userServicev1.UpdateUserResponse, error) {
	panic("implement me")
}
