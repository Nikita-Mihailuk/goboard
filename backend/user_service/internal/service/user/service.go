package user

import (
	"context"
	"github.com/Nikita-Mihailuk/goboard/backend/user_service/internal/domain/dto"
	"go.uber.org/zap"
)

type UserService struct {
	log          *zap.Logger
	userSaver    UserSaver
	userProvider UserProvider
	userUpdater  UserUpdater
}

func NewUserService(
	log *zap.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	userUpdater UserUpdater,
) *UserService {

	return &UserService{
		log:          log,
		userSaver:    userSaver,
		userProvider: userProvider,
		userUpdater:  userUpdater,
	}
}

type UserSaver interface {
	SaveUser(ctx context.Context, input dto.CreateUserInput) error
}

type UserProvider interface {
	GetUserByEmail(ctx context.Context, email string) (dto.LoginUserOutput, error)
	GetUserByID(ctx context.Context, id string) (dto.GetUserByIDOutput, error)
}

type UserUpdater interface {
	RefreshUser(ctx context.Context, id int64, input dto.UpdateUserInput) error
}
