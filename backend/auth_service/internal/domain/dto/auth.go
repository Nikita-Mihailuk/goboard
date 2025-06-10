package dto

type RegisterUserInput struct {
	Name     string
	Email    string
	Password string
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type LoginUserOutput struct {
	ID           int64
	PasswordHash string
	Role         string
}
