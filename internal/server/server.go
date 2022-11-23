package server

import (
	"net"

	"google.golang.org/grpc"

	"github.com/paramonies/ya-gophkeeper/internal/server/interceptor"
	"github.com/paramonies/ya-gophkeeper/internal/server/service"
	"github.com/paramonies/ya-gophkeeper/internal/store"
	pb "github.com/paramonies/ya-gophkeeper/pkg/gen/api/gophkeeper/v1"
	"github.com/paramonies/ya-gophkeeper/pkg/logger"
)

type GRPCServer struct {
	serv *grpc.Server
	log  *logger.Logger
	pb.UnimplementedUserServiceServer
	pb.UnimplementedPasswordServiceServer
}

func InitGRPCServer(con store.Connector, l *logger.Logger) (*GRPCServer, error) {
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.AuthCheckGRPC))

	userHandler := service.NewUserHandler(con, l)
	pb.RegisterUserServiceServer(grpcServer, userHandler)

	passwordHandler := service.NewPasswordHandler(con, l)
	pb.RegisterPasswordServiceServer(grpcServer, passwordHandler)

	return &GRPCServer{serv: grpcServer, log: l}, nil
}

func (s *GRPCServer) Start(address string) error {
	listen, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	s.log.Info("gRPC-server start")
	return s.serv.Serve(listen)
}

func (s *GRPCServer) ShutDown() error {
	s.serv.GracefulStop()
	return nil
}
