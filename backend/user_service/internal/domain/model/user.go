package model

type User struct {
	ID           int64
	Name         string
	Email        string
	PhotoURL     string
	PasswordHash string
	Role         string
}
