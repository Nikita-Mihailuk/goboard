package dto

import "mime/multipart"

type CreateUserInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserOutput struct {
	ID           int64
	PasswordHash string
	Role         string
}

type UpdateUserInput struct {
	ID              int64
	CurrentPassword string
	NewPassword     string
	Name            string
	FileHeader      *multipart.FileHeader
}

type GetUserByIDOutput struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	PhotoURL string `json:"photo_url"`
}
