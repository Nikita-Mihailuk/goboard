package auth_service

import "errors"

var (
	ErrUserExists         = errors.New("user already exists")
	ErrInternalGRPC       = errors.New("internal gRPC server error")
	ErrInvalidCredentials = errors.New("invalid credentials")
)
