package grpc

import (
	"fmt"
	userGRPC "github.com/Nikita-Mihailuk/goboard/backend/user_service/internal/delivery/grpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

type App struct {
	log        *zap.Logger
	grpcServer *grpc.Server
	port       string
}

func NewApp(log *zap.Logger, userService userGRPC.UserService, port string) *App {
	grpcServer := grpc.NewServer()
	userGRPC.RegisterGRPCServer(grpcServer, userService)
	return &App{log, grpcServer, port}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", a.port))
	if err != nil {
		return err
	}

	a.log.Info("starting grpc server", zap.String("address", lis.Addr().String()))

	if err = a.grpcServer.Serve(lis); err != nil {
		return err
	}
	return nil
}
func (a *App) Stop() {
	a.log.Info("stopping grpc server")
	a.grpcServer.GracefulStop()
}
