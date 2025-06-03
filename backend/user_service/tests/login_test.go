package tests

import (
	"testing"

	"github.com/Nikita-Mihailuk/goboard/backend/user_service/tests/suite"
	userServicev1 "github.com/Nikita-Mihailuk/protos_goboard/gen/go/user_service"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestLoginUser(t *testing.T) {
	ctx, s := suite.NewSuite(t)

	// создаём тестового пользователя
	email := gofakeit.Email()
	testUser := &userServicev1.CreateUserRequest{
		Email:    email,
		Password: "password",
		Name:     gofakeit.Name(),
	}
	_, err := s.UserClient.CreateUser(ctx, testUser)
	require.NoError(t, err)

	tests := []struct {
		name    string
		input   *userServicev1.LoginUserRequest
		wantErr error
	}{
		{
			name: "successful login",
			input: &userServicev1.LoginUserRequest{
				Email: email,
			},
			wantErr: nil,
		},
		{
			name: "empty email",
			input: &userServicev1.LoginUserRequest{
				Email: "",
			},
			wantErr: status.Error(codes.InvalidArgument, "missing fields"),
		},
		{
			name: "non-existent user",
			input: &userServicev1.LoginUserRequest{
				Email: gofakeit.Email(),
			},
			wantErr: status.Error(codes.NotFound, "user with email not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := s.UserClient.LoginUser(ctx, tt.input)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, status.Code(tt.wantErr), status.Code(err))
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, resp.UserId)
				assert.NotEmpty(t, resp.PasswordHash)
				assert.NotEmpty(t, resp.Role)
			}
		})
	}
}
