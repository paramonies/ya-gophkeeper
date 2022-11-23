package service

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	util "github.com/paramonies/ya-gophkeeper/internal/server/utils"
	"github.com/paramonies/ya-gophkeeper/internal/store"
	"github.com/paramonies/ya-gophkeeper/internal/store/dto"
	pb "github.com/paramonies/ya-gophkeeper/pkg/gen/api/gophkeeper/v1"
	"github.com/paramonies/ya-gophkeeper/pkg/logger"
)

type UserHandler struct {
	pb.UnimplementedUserServiceServer
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
	if req.GetLogin() == "" || req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "login and password must be specified")
	}

	h.log.Debug("input args", "login: ", req.GetLogin(), "password: ", req.GetPassword())

	resp, err := h.storage.Users().Register(ctx, &dto.RegisterRequest{
		Login:        req.GetLogin(),
		PasswordHash: util.EncryptPass(req.GetPassword()),
	})
	if err != nil {
		h.log.Error("error creating the user", err)
		return nil, handleError(err, "error creating the cluster")
	}

	token, err := util.JWTEncodeUserID(resp.UserID)
	if err != nil {
		h.log.Error("failed to create jwt token", err)
		return nil, status.Error(codes.Internal, "failed to create jwt token")
	}

	return &pb.RegisterUserResponse{
		UserID: resp.UserID,
		Jwt:    token,
	}, nil
}

// RegisterUser handler creates new user
func (h *UserHandler) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	h.log.Info("LoginUser handler")
	if req.GetLogin() == "" || req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "login and password must be specified")
	}

	h.log.Debug("input args", "login: ", req.GetLogin(), "password: ", req.GetPassword())

	resp, err := h.storage.Users().Login(ctx, &dto.LoginRequest{
		Login: req.GetLogin(),
	})
	if err != nil {
		h.log.Error("error getting the user", err)
		return nil, handleError(err, fmt.Sprintf("error getting the user %s", req.GetLogin()))
	}

	if util.EncryptPass(req.GetPassword()) != resp.PasswordHash {
		h.log.Error(fmt.Printf("access denied, wrong password for user %s", req.GetLogin()))
		return nil, status.Error(codes.Unauthenticated, "access denied, wrong password")
	}

	token, err := util.JWTEncodeUserID(resp.UserID)
	if err != nil {
		h.log.Error("failed to create jwt token", err)
		return nil, status.Error(codes.Internal, "failed to create jwt token")
	}

	return &pb.LoginUserResponse{
		UserID: resp.UserID,
		Jwt:    token,
	}, nil
}
