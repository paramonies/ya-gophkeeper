package client

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/paramonies/ya-gophkeeper/pkg/gen/api/gophkeeper/v1"
)

var clientConn *grpc.ClientConn

type ClientSet struct {
	UserClient pb.UserServiceClient
	PwdClient  pb.PasswordServiceClient
}

func CreateClientSet(serverPath string) (*ClientSet, error) {
	conn, err := grpc.Dial(serverPath, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	clientConn = conn

	return &ClientSet{
		UserClient: pb.NewUserServiceClient(conn),
		PwdClient:  pb.NewPasswordServiceClient(conn),
	}, nil
}

func ConnDown() error {
	return clientConn.Close()
}
