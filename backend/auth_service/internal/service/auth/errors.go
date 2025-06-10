package auth

import "errors"

var (
	ErrUserExists          = errors.New("user already exists")
	ErrInternalUserService = errors.New("internal user service error")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
)
