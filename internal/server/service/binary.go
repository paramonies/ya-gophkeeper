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

type BinaryHandler struct {
	pb.UnimplementedBinaryServiceServer
	storage store.Connector
	log     *logger.Logger
}

func NewBinaryHandler(s store.Connector, l *logger.Logger) *BinaryHandler {
	return &BinaryHandler{
		storage: s,
		log:     l,
	}
}

func (h *BinaryHandler) CreateBinary(ctx context.Context, req *pb.CreateBinaryRequest) (*pb.CreateBinaryResponse, error) {
	h.log.Debug("CreateBinary handler")
	if req.GetVersion() < 1 || req.GetTitle() == "" || req.GetData() == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid arguments")
	}

	userID := utils.GetUserIDFromCTX(ctx)

	res, err := h.storage.Binaries().GetByTitle(ctx, &dto.GetBinaryByTitleRequest{
		Title:  req.GetTitle(),
		UserID: *userID,
	})
	if err != nil && !core.IsNotFound(err) {
		msg := fmt.Sprintf("error getting the binary for %s from db", req.GetTitle())
		h.log.Error(msg, err)
		return nil, handleError(err, msg)
	}

	if err != nil && core.IsNotFound(err) {
		_, errCr := h.storage.Binaries().Create(ctx, &dto.CreateBinaryRequest{
			UserID:  *userID,
			Title:   req.GetTitle(),
			Data:    req.GetData(),
			Meta:    req.GetMeta(),
			Version: req.GetVersion(),
		})
		if errCr != nil {
			msg := fmt.Sprintf("error creating the binary for title %s", req.GetTitle())
			h.log.Error(msg, errCr)
			return nil, handleError(errCr, msg)
		}

		return &pb.CreateBinaryResponse{
			Status: "New binary created",
		}, nil
	}

	if req.GetVersion() <= res.Version {
		return nil, status.Error(codes.AlreadyExists, newerVersionDetected)
	}

	_, errCr := h.storage.Binaries().Create(ctx, &dto.CreateBinaryRequest{
		UserID:  *userID,
		Title:   req.GetTitle(),
		Data:    req.GetData(),
		Meta:    req.GetMeta(),
		Version: req.GetVersion(),
	})
	if errCr != nil {
		msg := fmt.Sprintf("error creating the binary for title %s", req.GetTitle())
		h.log.Error(msg, err)
		return nil, handleError(err, msg)
	}

	return &pb.CreateBinaryResponse{
		Status: "New binary created",
	}, nil
}

func (h *BinaryHandler) GetBinary(ctx context.Context, req *pb.GetBinaryRequest) (*pb.GetBinaryResponse, error) {
	h.log.Debug("GetBinary handler")
	if req.GetTitle() == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}

	userID := utils.GetUserIDFromCTX(ctx)

	res, err := h.storage.Binaries().GetByTitle(ctx, &dto.GetBinaryByTitleRequest{
		Title:  req.GetTitle(),
		UserID: *userID,
	})
	if err != nil {
		if core.IsNotFound(err) {
			msg := fmt.Sprintf("data for title %s not found", req.GetTitle())
			h.log.Error(msg, err)
			return nil, handleError(err, msg)
		}

		msg := fmt.Sprintf("failed to obtain latest data for title %s", req.GetTitle())
		h.log.Error(msg, err)
		return nil, handleError(err, msg)
	}

	return &pb.GetBinaryResponse{
		Title:   req.GetTitle(),
		Data:    res.Data,
		Meta:    res.Meta,
		Version: res.Version,
	}, nil
}

func (h *BinaryHandler) DeleteBinary(ctx context.Context, req *pb.DeleteBinaryRequest) (*pb.DeleteBinaryResponse, error) {
	h.log.Debug("DeleteBinary handler")
	if req.GetTitle() == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}

	userID := utils.GetUserIDFromCTX(ctx)

	err := h.storage.Binaries().Delete(ctx, &dto.DeleteBinaryRequest{
		Title:  req.GetTitle(),
		UserID: *userID,
	})
	if err != nil {
		msg := fmt.Sprintf("failed to delete data for title %s", req.GetTitle())
		h.log.Error(msg, err)
		return nil, status.Error(codes.Internal, msg)
	}

	return &pb.DeleteBinaryResponse{
		Status: "success",
	}, nil

}
