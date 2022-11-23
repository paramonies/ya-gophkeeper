package service

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/paramonies/ya-gophkeeper/internal/core"
	"github.com/paramonies/ya-gophkeeper/internal/server/utils"
	"github.com/paramonies/ya-gophkeeper/internal/store"
	"github.com/paramonies/ya-gophkeeper/internal/store/dto"
	pb "github.com/paramonies/ya-gophkeeper/pkg/gen/api/gophkeeper/v1"
	"github.com/paramonies/ya-gophkeeper/pkg/logger"
)

var (
	newerVersionDetected = "newer version found in database. please synchronize you app to get the most actual data."
)

type PasswordHandler struct {
	pb.UnimplementedPasswordServiceServer
	storage store.Connector
	log     *logger.Logger
}

func NewPasswordHandler(s store.Connector, l *logger.Logger) *PasswordHandler {
	return &PasswordHandler{
		storage: s,
		log:     l,
	}
}

func (h *PasswordHandler) CreatePassword(ctx context.Context, req *pb.CreatePasswordRequest) (*pb.CreatePasswordResponse, error) {
	h.log.Info("CreatePassword handler")
	if req.GetVersion() < 1 || req.GetLogin() == "" || req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid arguments")
	}

	userID := utils.GetUserIDFromCTX(ctx)

	res, err := h.storage.Passwords().GetByID(ctx, &dto.GetByIDRequest{
		Login:  req.GetLogin(),
		UserID: *userID,
	})
	if err != nil && !core.IsNotFound(err) {
		h.log.Error(fmt.Sprintf("error getting the password for %s from db", req.GetLogin()), err)
		return nil, handleError(err, fmt.Sprintf("error getting the password for %s from db", req.GetLogin()))
	}

	if err != nil && core.IsNotFound(err) {
		_, errCr := h.storage.Passwords().Create(ctx, &dto.CreateRequest{
			UserID:   *userID,
			Login:    req.GetLogin(),
			Password: req.GetPassword(),
			Meta:     req.GetMeta(),
			Version:  req.GetVersion(),
		})
		if errCr != nil {
			h.log.Error(fmt.Sprintf("error creating the password for login %s", req.Login), err)
			return nil, handleError(err, fmt.Sprintf("error creating the password for login %s", req.Login))
		}

		return &pb.CreatePasswordResponse{
			Status: "New password created",
		}, nil
	}

	if req.GetVersion() <= res.Version {
		return nil, status.Error(codes.AlreadyExists, newerVersionDetected)
	}

	_, errCr := h.storage.Passwords().Create(ctx, &dto.CreateRequest{
		UserID:   *userID,
		Login:    req.GetLogin(),
		Password: req.GetPassword(),
		Meta:     req.GetMeta(),
		Version:  req.GetVersion(),
	})
	if errCr != nil {
		h.log.Error(fmt.Sprintf("error creating the password for login %s", req.Login), err)
		return nil, handleError(err, fmt.Sprintf("error creating the password for login %s", req.Login))
	}

	return &pb.CreatePasswordResponse{
		Status: "New password created",
	}, nil
}
