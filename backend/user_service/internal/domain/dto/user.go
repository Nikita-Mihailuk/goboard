package dto

import "database/sql"

type CreateUserInput struct {
	Name         string
	Email        string
	PasswordHash string
}

type LoginUserOutput struct {
	ID           int64
	PasswordHash string
	Role         string
}

type UpdateUserInput struct {
	ID              int64
	CurrentPassword string
	InputPassword   string
	Name            string
	PhotoUrl        string
}

type UserForUpdate struct {
	ID           int64
	PasswordHash string
	Name         string
	PhotoUrl     sql.NullString
}

type GetUserByIDOutput struct {
	Email    string
	Name     string
	PhotoURL sql.NullString
}

type UpdateUserMessage struct {
	UserID       int64  `json:"user_id"`
	UserName     string `json:"user_name"`
	UserPhotoURL string `json:"user_photo_url"`
}
