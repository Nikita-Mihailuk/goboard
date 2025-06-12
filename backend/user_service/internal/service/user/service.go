package user

import (
	"context"
	"github.com/Nikita-Mihailuk/goboard/backend/user_service/internal/domain/dto"
	"github.com/Nikita-Mihailuk/goboard/backend/user_service/internal/infrastructure/kafka"
	"go.uber.org/zap"
)

type UserService struct {
	log           *zap.Logger
	userSaver     UserSaver
	userProvider  UserProvider
	userUpdater   UserUpdater
	kafkaProducer *kafka.Producer
	producerTopic string
}

func NewUserService(
	log *zap.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	userUpdater UserUpdater,
	kafkaProducer *kafka.Producer,
	producerTopic string,
) *UserService {

	return &UserService{
		log:           log,
		userSaver:     userSaver,
		userProvider:  userProvider,
		userUpdater:   userUpdater,
		kafkaProducer: kafkaProducer,
		producerTopic: producerTopic,
	}
}

type UserSaver interface {
	SaveUser(ctx context.Context, input dto.CreateUserInput) error
}

type UserProvider interface {
	GetUserByEmail(ctx context.Context, email string) (dto.LoginUserOutput, error)
	GetUserByID(ctx context.Context, id int64) (dto.GetUserByIDOutput, error)
	GetUserUpdateByID(ctx context.Context, id int64) (dto.UserForUpdate, error)
}

type UserUpdater interface {
	RefreshUser(ctx context.Context, input dto.UserForUpdate) error
}
