package suite

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Nikita-Mihailuk/goboard/backend/user_service/internal/config"
	userServicev1 "github.com/Nikita-Mihailuk/protos_goboard/gen/go/user_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Suite struct {
	*testing.T
	Cfg        *config.Config
	UserClient userServicev1.UserClient
}

const gRPCTimeout time.Duration = time.Second * 5

func NewSuite(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.LoadConfigByPath("../config/local_test.yml")
	ctx, cancel := context.WithTimeout(context.Background(), gRPCTimeout)

	t.Cleanup(func() {
		t.Helper()
		cancel()
	})

	cc, err := grpc.DialContext(
		context.Background(),
		fmt.Sprintf("localhost:%s", cfg.GRPCServer.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		panic(err)
	}

	return ctx, &Suite{
		T:          t,
		Cfg:        cfg,
		UserClient: userServicev1.NewUserClient(cc),
	}
}
