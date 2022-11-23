package grpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/paramonies/ya-gophkeeper/pkg/gen/api/gophkeeper/v1"
)

var clientConn *grpc.ClientConn

// DialUp initiates a connection between the client and the server. Address taken from cfg.GrpcServerPath.
func DialUpUser(serverPath string) (pb.UserServiceClient, error) {
	conn, err := grpc.Dial(serverPath, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	clientConn = conn
	return pb.NewUserServiceClient(conn), nil
}

func DialUpPass(serverPath string) (pb.PasswordServiceClient, error) {
	conn, err := grpc.Dial(serverPath, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	clientConn = conn
	return pb.NewPasswordServiceClient(conn), nil
}

func ConnDown() error {
	return clientConn.Close()
}

func ActiveConnection() bool {
	return clientConn != nil
}
