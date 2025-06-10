package user

import (
	"context"
	"errors"
	"github.com/Nikita-Mihailuk/goboard/backend/user_service/internal/domain/dto"
	"github.com/Nikita-Mihailuk/goboard/backend/user_service/internal/storage/postgres"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func (s *UserService) CreateUser(ctx context.Context, input dto.CreateUserInput) error {
	err := s.userSaver.SaveUser(ctx, input)
	if err != nil {
		if errors.Is(err, postgres.ErrUserExists) {
			s.log.Error("user already exists", zap.Error(err))
			return ErrUserExists
		}
		s.log.Error("failed create user", zap.Error(err))
		return err
	}

	s.log.Info("created user", zap.String("email", input.Email))
	return nil
}

func (s *UserService) GetLoginUser(ctx context.Context, email string) (dto.LoginUserOutput, error) {
	outputUser, err := s.userProvider.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, postgres.ErrUserNotFound) {
			s.log.Error("user not found", zap.Error(err))
			return dto.LoginUserOutput{}, ErrUserNotFound
		}
		s.log.Error("failed get user by email", zap.Error(err))
		return dto.LoginUserOutput{}, err
	}

	s.log.Info("get user by email", zap.String("email", email))
	return outputUser, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id int64) (dto.GetUserByIDOutput, error) {
	outputUser, err := s.userProvider.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, postgres.ErrUserNotFound) {
			s.log.Error("user not found", zap.Error(err))
			return dto.GetUserByIDOutput{}, ErrUserNotFound
		}
		s.log.Error("failed get user by id", zap.Error(err))
		return dto.GetUserByIDOutput{}, err
	}

	s.log.Info("get user by id", zap.Int64("id", id))
	return outputUser, nil
}

func (s *UserService) UpdateUser(ctx context.Context, input dto.UpdateUserInput) error {
	user, err := s.userProvider.GetUserUpdateByID(ctx, input.ID)
	if err != nil {
		if errors.Is(err, postgres.ErrUserNotFound) {
			s.log.Error("user not found", zap.Error(err))
			return ErrUserNotFound
		}
		s.log.Error("failed get password by id", zap.Error(err))
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.CurrentPassword))
	if err != nil {
		s.log.Error("failed compare current password", zap.Error(err))
		return ErrInvalidPassword
	}

	if input.InputPassword != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.InputPassword), bcrypt.DefaultCost)
		if err != nil {
			s.log.Error("failed to hash password", zap.Error(err))
			return err
		}
		user.PasswordHash = string(hashedPassword)
	}

	if input.Name != "" {
		user.Name = input.Name
	}

	if input.PhotoUrl != "" {
		user.PhotoUrl.String = input.PhotoUrl
	}

	err = s.userUpdater.RefreshUser(ctx, user)
	if err != nil {
		s.log.Error("failed update user", zap.Error(err))
		return err
	}

	s.log.Info("update user", zap.Int64("id", user.ID))
	return nil
}
