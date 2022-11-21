package interceptor

import (
	"context"
	"log"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/paramonies/ya-gophkeeper/internal/server/utils"
)

var (
	SkipCheckMethods = map[string]struct{}{
		"/proto.GophkeeperService/RegisterUser": {},
		"/proto.GophkeeperService/LoginUser":    {},
	}
)

// AuthCheckGRPC interceptor verifies the authentication bearer token.
func AuthCheckGRPC(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	log.Println("--> unary interceptor: ", info.FullMethod)
	_, ok := SkipCheckMethods[info.FullMethod]
	if ok {
		return handler(ctx, req)
	}

	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}

	userID, err := utils.JWTDecodeUserID(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}

	return handler(utils.SetUserIDToCTX(ctx, userID), req)
}
