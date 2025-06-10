package grpc

import (
	"context"
	"errors"
	"github.com/Nikita-Mihailuk/goboard/backend/user_service/internal/domain/dto"
	"github.com/Nikita-Mihailuk/goboard/backend/user_service/internal/service/user"
	userServicev1 "github.com/Nikita-Mihailuk/protos_goboard/gen/go/user_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *serverGRPC) CreateUser(ctx context.Context, req *userServicev1.CreateUserRequest) (*userServicev1.CreateUserResponse, error) {
	inputUser, err := validateCreateUserRequest(req)
	if err != nil {
		return nil, err
	}
	err = s.userService.CreateUser(ctx, inputUser)
	if err != nil {
		if errors.Is(err, user.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &userServicev1.CreateUserResponse{}, nil
}

func (s *serverGRPC) LoginUser(ctx context.Context, req *userServicev1.LoginUserRequest) (*userServicev1.LoginUserResponse, error) {
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "missing fields")
	}

	outputUser, err := s.userService.GetLoginUser(ctx, req.GetEmail())
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user with email not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &userServicev1.LoginUserResponse{
		UserId:       outputUser.ID,
		PasswordHash: outputUser.PasswordHash,
		Role:         outputUser.Role,
	}, nil
}

func (s *serverGRPC) GetUserByID(ctx context.Context, req *userServicev1.GetUserByIDRequest) (*userServicev1.GetUserByIDResponse, error) {
	if req.GetUserId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "missing fields")
	}

	outputUser, err := s.userService.GetUserByID(ctx, req.GetUserId())
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user with ID not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &userServicev1.GetUserByIDResponse{
		Email:    outputUser.Email,
		Name:     outputUser.Name,
		PhotoUrl: outputUser.PhotoURL.String,
	}, nil
}

func (s *serverGRPC) UpdateUser(ctx context.Context, req *userServicev1.UpdateUserRequest) (*userServicev1.UpdateUserResponse, error) {
	inputUser, err := validateUpdateUserRequest(req)
	if err != nil {
		return nil, err
	}

	err = s.userService.UpdateUser(ctx, inputUser)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user with ID not found")
		}
		if errors.Is(err, user.ErrInvalidPassword) {
			return nil, status.Error(codes.InvalidArgument, "invalid password")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &userServicev1.UpdateUserResponse{}, nil
}

func validateCreateUserRequest(req *userServicev1.CreateUserRequest) (dto.CreateUserInput, error) {
	if req.GetEmail() == "" || req.GetPassword() == "" || req.GetName() == "" {
		return dto.CreateUserInput{}, status.Errorf(codes.InvalidArgument, "missing fields")
	}

	return dto.CreateUserInput{
		Email:        req.GetEmail(),
		PasswordHash: req.GetPassword(),
		Name:         req.GetName(),
	}, nil
}

func validateUpdateUserRequest(req *userServicev1.UpdateUserRequest) (dto.UpdateUserInput, error) {
	if req.GetCurrentPassword() == "" || req.GetUserId() == 0 {
		return dto.UpdateUserInput{}, status.Errorf(codes.InvalidArgument, "missing fields")
	}

	return dto.UpdateUserInput{
		ID:              req.GetUserId(),
		CurrentPassword: req.GetCurrentPassword(),
		InputPassword:   req.GetNewPassword(),
		Name:            req.GetName(),
		PhotoUrl:        req.GetPhotoUrl(),
	}, nil
}
