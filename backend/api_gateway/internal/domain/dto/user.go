package dto

import "mime/multipart"

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
	PhotoURL string `json:"photo_url,omitempty"`
}
