package postgres

import (
	"context"
	"github.com/Nikita-Mihailuk/goboard/backend/user_service/internal/domain/dto"
)

func (s *Storage) SaveUser(ctx context.Context, input dto.CreateUserInput) error {
	panic("implement me")
}

func (s *Storage) GetUserByEmail(ctx context.Context, email string) (dto.LoginUserOutput, error) {
	panic("implement me")
}

func (s *Storage) GetUserByID(ctx context.Context, id string) (dto.GetUserByIDOutput, error) {
	panic("implement me")
}

func (s *Storage) RefreshUser(ctx context.Context, id int64, input dto.UpdateUserInput) error {
	panic("implement me")
}
