package user_service

import "errors"

var (
	ErrUserExists         = errors.New("user already exists")
	ErrInternalGRPC       = errors.New("internal gRPC server error")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidPassword    = errors.New("invalid password")
)
