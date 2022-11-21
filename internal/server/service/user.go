package service

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	util "github.com/paramonies/ya-gophkeeper/internal/server/utils"
	"github.com/paramonies/ya-gophkeeper/internal/store"
	"github.com/paramonies/ya-gophkeeper/internal/store/dto"
	pb "github.com/paramonies/ya-gophkeeper/pkg/gen/api/gophkeeper/v1"
	"github.com/paramonies/ya-gophkeeper/pkg/logger"
)

type UserHandler struct {
	pb.UnimplementedGophkeeperServiceServer
	storage store.Connector
	log     *logger.Logger
}

func NewUserHandler(s store.Connector, l *logger.Logger) *UserHandler {
	return &UserHandler{
		storage: s,
		log:     l,
	}
}

// RegisterUser handler creates new user
func (h *UserHandler) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	h.log.Info("RegisterUser handler")
	if req.ServiceLogin == "" || req.ServicePass == "" {
		return nil, status.Error(codes.InvalidArgument, "login and password must be specified")
	}

	h.log.Debug("input args", "login: ", req.GetServiceLogin(), "password: ", req.GetServicePass())

	resp, err := h.storage.Users().Register(ctx, &dto.RegisterRequest{
		Login:        req.GetServiceLogin(),
		PasswordHash: util.EncryptPass(req.GetServicePass()),
	})
	if err != nil {
		h.log.Error("error creating the user", err)
		return nil, handleError(err, "error creating the cluster")
	}

	userJWT, err := util.JWTEncodeUserID(resp.UserID)
	if err != nil {
		h.log.Error("failed to create jwt token", err)
		return nil, status.Error(codes.Internal, "failed to create jwt token")
	}

	return &pb.RegisterUserResponse{
		UserID: resp.UserID,
		Jwt:    userJWT,
	}, nil
}
