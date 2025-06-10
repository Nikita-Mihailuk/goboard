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
