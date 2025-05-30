package user

import (
	"context"
	"github.com/Nikita-Mihailuk/goboard/backend/user_service/internal/domain/dto"
)

func (s *UserService) CreateUser(ctx context.Context, input dto.CreateUserInput) error {
	panic("implement me")
}

func (s *UserService) LoginUser(ctx context.Context, email string) (dto.LoginUserOutput, error) {
	panic("implement me")
}

func (s *UserService) GetUserByID(ctx context.Context, id int64) (dto.GetUserByIDOutput, error) {
	panic("implement me")
}

func (s *UserService) UpdateUser(ctx context.Context, id int64, input dto.UpdateUserInput) error {
	panic("implement me")
}
