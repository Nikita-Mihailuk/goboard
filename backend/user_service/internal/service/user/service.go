package user

import (
	"go.uber.org/zap"
)

type UserService struct {
	log *zap.Logger
}

func NewEmployeeService(
	log *zap.Logger,
) *UserService {

	return &UserService{
		log: log,
	}
}
