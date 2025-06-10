package user_service

import (
	"context"
	"fmt"
	"github.com/Nikita-Mihailuk/goboard/backend/api_gateway/internal/domain/dto"
	userServicev1 "github.com/Nikita-Mihailuk/protos_goboard/gen/go/user_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"mime/multipart"
	"os"
)

type UserClient struct {
	api userServicev1.UserClient
}

func NewUserClient(ctx context.Context, addr string) (*UserClient, error) {
	cc, err := grpc.DialContext(ctx, addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &UserClient{api: userServicev1.NewUserClient(cc)}, nil
}

func (c *UserClient) GetUserByID(ctx context.Context, id int64) (dto.GetUserByIDOutput, error) {
	resp, err := c.api.GetUserByID(ctx, &userServicev1.GetUserByIDRequest{
		UserId: id,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				return dto.GetUserByIDOutput{}, ErrUserNotFound
			default:
				return dto.GetUserByIDOutput{}, ErrInternalGRPC
			}
		}
		return dto.GetUserByIDOutput{}, err
	}

	return dto.GetUserByIDOutput{
		Email:    resp.GetEmail(),
		Name:     resp.GetName(),
		PhotoURL: resp.GetPhotoUrl(),
	}, nil
}

func (c *UserClient) UpdateUser(ctx context.Context, input dto.UpdateUserInput) error {
	// TODO: redesign the photo saving logic on s3
	var photoURL string
	if input.FileHeader != nil {
		filePath := fmt.Sprintf("static/%d_%s", input.ID, input.FileHeader.Filename)
		err := SaveFile(input.FileHeader, filePath)
		if err != nil {
			return err
		}
		photoURL = filePath
	}
	_, err := c.api.UpdateUser(ctx, &userServicev1.UpdateUserRequest{
		UserId:          input.ID,
		CurrentPassword: input.CurrentPassword,
		NewPassword:     input.NewPassword,
		PhotoUrl:        photoURL,
		Name:            input.Name,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				return ErrUserNotFound
			case codes.InvalidArgument:
				return ErrInvalidPassword
			default:
				return ErrInternalGRPC
			}
		}
		return err
	}

	return nil
}

func SaveFile(fileHeader *multipart.FileHeader, path string) error {
	src, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}
	return nil
}
