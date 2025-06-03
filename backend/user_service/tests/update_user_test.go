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

func TestUpdateUser(t *testing.T) {
	ctx, s := suite.NewSuite(t)

	// создаём тестового пользователя
	testUser := &userServicev1.CreateUserRequest{
		Email:    gofakeit.Email(),
		Password: "password",
		Name:     gofakeit.Name(),
	}
	_, err := s.UserClient.CreateUser(ctx, testUser)
	require.NoError(t, err)

	// получаем id созданного пользователя
	loginResp, err := s.UserClient.LoginUser(ctx, &userServicev1.LoginUserRequest{
		Email: testUser.Email,
	})
	require.NoError(t, err)
	testUserID := loginResp.UserId

	tests := []struct {
		name    string
		input   *userServicev1.UpdateUserRequest
		wantErr error
	}{
		{
			name: "successful name update",
			input: &userServicev1.UpdateUserRequest{
				UserId:          testUserID,
				CurrentPassword: "password",
				Name:            gofakeit.Name(),
			},
			wantErr: nil,
		},
		{
			name: "successful password update",
			input: &userServicev1.UpdateUserRequest{
				UserId:          testUserID,
				CurrentPassword: "password",
				NewPassword:     "newpassword",
			},
			wantErr: nil,
		},
		{
			name: "successful photo update",
			input: &userServicev1.UpdateUserRequest{
				UserId:          testUserID,
				CurrentPassword: "newpassword",
				PhotoUrl:        gofakeit.URL(),
			},
			wantErr: nil,
		},
		{
			name: "invalid current password",
			input: &userServicev1.UpdateUserRequest{
				UserId:          testUserID,
				CurrentPassword: "wrongpassword",
				Name:            gofakeit.Name(),
			},
			wantErr: status.Error(codes.InvalidArgument, "invalid password"),
		},
		{
			name: "user not found",
			input: &userServicev1.UpdateUserRequest{
				UserId:          99999,
				CurrentPassword: "password",
				Name:            gofakeit.Name(),
			},
			wantErr: status.Error(codes.NotFound, "user with ID not found"),
		},
		{
			name: "missing user ID",
			input: &userServicev1.UpdateUserRequest{
				UserId:          0,
				CurrentPassword: "password",
				Name:            gofakeit.Name(),
			},
			wantErr: status.Error(codes.InvalidArgument, "missing fields"),
		},
		{
			name: "missing current password",
			input: &userServicev1.UpdateUserRequest{
				UserId:          testUserID,
				CurrentPassword: "",
				Name:            gofakeit.Name(),
			},
			wantErr: status.Error(codes.InvalidArgument, "missing fields"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := s.UserClient.UpdateUser(ctx, tt.input)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, status.Code(tt.wantErr), status.Code(err))
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				require.NoError(t, err)

				// проверяем сохранение изменений
				resp, err := s.UserClient.GetUserByID(ctx, &userServicev1.GetUserByIDRequest{
					UserId: tt.input.UserId,
				})
				require.NoError(t, err)

				if tt.input.Name != "" {
					assert.Equal(t, tt.input.Name, resp.Name)
				}
				if tt.input.PhotoUrl != "" {
					assert.Equal(t, tt.input.PhotoUrl, resp.PhotoUrl)
				}
			}
		})
	}
}
