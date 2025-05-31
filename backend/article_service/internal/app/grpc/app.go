package grpc

import (
	"fmt"
	articleGRPC "github.com/Nikita-Mihailuk/goboard/backend/article_service/internal/delivery/grpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

type App struct {
	log        *zap.Logger
	grpcServer *grpc.Server
	port       string
}

func NewApp(log *zap.Logger, articleService articleGRPC.ArticleService, port string) *App {
	grpcServer := grpc.NewServer()
	articleGRPC.RegisterGRPCServer(grpcServer, articleService)
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
