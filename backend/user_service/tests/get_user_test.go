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

func TestGetUserByID(t *testing.T) {
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
		input   *userServicev1.GetUserByIDRequest
		want    *userServicev1.GetUserByIDResponse
		wantErr error
	}{
		{
			name: "successful get user",
			input: &userServicev1.GetUserByIDRequest{
				UserId: testUserID,
			},
			want: &userServicev1.GetUserByIDResponse{
				Email: testUser.Email,
				Name:  testUser.Name,
			},
			wantErr: nil,
		},
		{
			name: "user not found",
			input: &userServicev1.GetUserByIDRequest{
				UserId: 99999,
			},
			want:    nil,
			wantErr: status.Error(codes.NotFound, "user with ID not found"),
		},
		{
			name: "invalid ID",
			input: &userServicev1.GetUserByIDRequest{
				UserId: 0,
			},
			want:    nil,
			wantErr: status.Error(codes.InvalidArgument, "missing fields"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := s.UserClient.GetUserByID(ctx, tt.input)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, status.Code(tt.wantErr), status.Code(err))
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want.Email, resp.Email)
				assert.Equal(t, tt.want.Name, resp.Name)
			}
		})
	}
}
