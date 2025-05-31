package grpc

import (
	"context"
	"github.com/Nikita-Mihailuk/goboard/backend/user_service/internal/domain/dto"
	userServicev1 "github.com/Nikita-Mihailuk/protos_goboard/gen/go/user_service"
	"google.golang.org/grpc"
)

type UserService interface {
	CreateUser(ctx context.Context, input dto.CreateUserInput) error
	GetLoginUser(ctx context.Context, email string) (dto.LoginUserOutput, error)
	GetUserByID(ctx context.Context, id int64) (dto.GetUserByIDOutput, error)
	UpdateUser(ctx context.Context, input dto.UpdateUserInput) error
}

type serverGRPC struct {
	userService UserService
	userServicev1.UnimplementedUserServer
}

func RegisterGRPCServer(grpcServer *grpc.Server, userService UserService) {
	userServicev1.RegisterUserServer(grpcServer, &serverGRPC{userService: userService})
}
