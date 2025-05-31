package user_service

import (
	"context"
	userServicev1 "github.com/Nikita-Mihailuk/protos_goboard/gen/go/user_service"
	"google.golang.org/grpc"
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
