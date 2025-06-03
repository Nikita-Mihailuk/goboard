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

func TestCreateUser(t *testing.T) {
	ctx, s := suite.NewSuite(t)

	// создаем тестового пользователя для проверки дубликата
	existingEmail := gofakeit.Email()
	existingUser := &userServicev1.CreateUserRequest{
		Email:    existingEmail,
		Password: "password",
		Name:     gofakeit.Name(),
	}
	_, err := s.UserClient.CreateUser(ctx, existingUser)
	require.NoError(t, err)

	tests := []struct {
		name    string
		input   *userServicev1.CreateUserRequest
		wantErr error
	}{
		{
			name: "successful user creation",
			input: &userServicev1.CreateUserRequest{
				Email:    gofakeit.Email(),
				Password: gofakeit.Password(true, true, true, true, false, 10),
				Name:     gofakeit.Name(),
			},
			wantErr: nil,
		},
		{
			name: "empty email",
			input: &userServicev1.CreateUserRequest{
				Email:    "",
				Password: gofakeit.Password(true, true, true, true, false, 10),
				Name:     gofakeit.Name(),
			},
			wantErr: status.Error(codes.InvalidArgument, "missing fields"),
		},
		{
			name: "empty password",
			input: &userServicev1.CreateUserRequest{
				Email:    gofakeit.Email(),
				Password: "",
				Name:     gofakeit.Name(),
			},
			wantErr: status.Error(codes.InvalidArgument, "missing fields"),
		},
		{
			name: "empty name",
			input: &userServicev1.CreateUserRequest{
				Email:    gofakeit.Email(),
				Password: gofakeit.Password(true, true, true, true, false, 10),
				Name:     "",
			},
			wantErr: status.Error(codes.InvalidArgument, "missing fields"),
		},
		{
			name: "user already exists",
			input: &userServicev1.CreateUserRequest{
				Email:    existingEmail,
				Password: gofakeit.Password(true, true, true, true, false, 10),
				Name:     gofakeit.Name(),
			},
			wantErr: status.Error(codes.AlreadyExists, "user already exists"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := s.UserClient.CreateUser(ctx, tt.input)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, status.Code(tt.wantErr), status.Code(err))
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				require.NoError(t, err)

				// проверяем создался ли пользователь
				loginResp, err := s.UserClient.LoginUser(ctx, &userServicev1.LoginUserRequest{
					Email: tt.input.Email,
				})
				require.NoError(t, err)
				assert.NotEmpty(t, loginResp.UserId)
			}
		})
	}
}
