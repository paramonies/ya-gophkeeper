package server

import (
	"net"

	"github.com/paramonies/ya-gophkeeper/pkg/graceful"

	"google.golang.org/grpc"

	"github.com/paramonies/ya-gophkeeper/internal/server/interceptor"
	"github.com/paramonies/ya-gophkeeper/internal/server/service"
	"github.com/paramonies/ya-gophkeeper/internal/store"
	pb "github.com/paramonies/ya-gophkeeper/pkg/gen/api/gophkeeper/v1"
	"github.com/paramonies/ya-gophkeeper/pkg/logger"
)

// RunGRPCServer runs grpc api server.
func RunGRPCServer(addr string, con store.Connector, l *logger.Logger) error {
	l.Info("start listening API server", "address", addr)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.AuthCheckGRPC),
	)

	graceful.AddCallback(func() error {
		grpcServer.GracefulStop()
		return nil
	})

	userHandler := service.NewUserHandler(con, l)
	passwordHandler := service.NewPasswordHandler(con, l)
	textHandler := service.NewTextHandler(con, l)
	binaryHandler := service.NewBinaryHandler(con, l)

	pb.RegisterUserServiceServer(grpcServer, userHandler)
	pb.RegisterPasswordServiceServer(grpcServer, passwordHandler)
	pb.RegisterTextServiceServer(grpcServer, textHandler)
	pb.RegisterBinaryServiceServer(grpcServer, binaryHandler)

	go func() {
		listenErr := grpcServer.Serve(listener)
		if listenErr != nil {
			l.Error("failed to serve gRPC API server", listenErr)
			graceful.ShutdownNow()
		}
	}()

	return nil
}
