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
	TextClient pb.TextServiceClient
	BinClient  pb.BinaryServiceClient
	CardClient pb.CardServiceClient
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
		TextClient: pb.NewTextServiceClient(conn),
		BinClient:  pb.NewBinaryServiceClient(conn),
		CardClient: pb.NewCardServiceClient(conn),
	}, nil
}

func ConnDown() error {
	return clientConn.Close()
}
