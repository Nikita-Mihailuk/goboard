package dto

type CreateUserInput struct {
	Name     string
	Email    string
	Password string
}

type LoginUserOutput struct {
	ID           int64
	PasswordHash string
	Role         string
}

type UpdateUserInput struct {
	CurrentPassword string
	InputPassword   string
	Name            string
	PhotoUrl        string
}

type GetUserByIDOutput struct {
	Email    string
	Name     string
	PhotoURL string
}
