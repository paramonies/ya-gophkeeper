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

type TextHandler struct {
	pb.UnimplementedTextServiceServer
	storage store.Connector
	log     *logger.Logger
}

func NewTextHandler(s store.Connector, l *logger.Logger) *TextHandler {
	return &TextHandler{
		storage: s,
		log:     l,
	}
}

func (h *TextHandler) CreateText(ctx context.Context, req *pb.CreateTextRequest) (*pb.CreateTextResponse, error) {
	h.log.Debug("CreateText handler")
	if req.GetVersion() < 1 || req.GetTitle() == "" || req.GetData() == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid arguments")
	}

	userID := utils.GetUserIDFromCTX(ctx)

	res, err := h.storage.Texts().GetByTitle(ctx, &dto.GetTextByTitleRequest{
		Title:  req.GetTitle(),
		UserID: *userID,
	})
	if err != nil && !core.IsNotFound(err) {
		msg := fmt.Sprintf("error getting the text for %s from db", req.GetTitle())
		h.log.Error(msg, err)
		return nil, handleError(err, msg)
	}

	if err != nil && core.IsNotFound(err) {
		_, errCr := h.storage.Texts().Create(ctx, &dto.CreateTextRequest{
			UserID:  *userID,
			Title:   req.GetTitle(),
			Data:    req.GetData(),
			Meta:    req.GetMeta(),
			Version: req.GetVersion(),
		})
		if errCr != nil {
			msg := fmt.Sprintf("error creating the text for title %s", req.GetTitle())
			h.log.Error(msg, errCr)
			return nil, handleError(errCr, msg)
		}

		return &pb.CreateTextResponse{
			Status: "New text created",
		}, nil
	}

	if req.GetVersion() <= res.Version {
		return nil, status.Error(codes.AlreadyExists, newerVersionDetected)
	}

	_, errCr := h.storage.Texts().Create(ctx, &dto.CreateTextRequest{
		UserID:  *userID,
		Title:   req.GetTitle(),
		Data:    req.GetData(),
		Meta:    req.GetMeta(),
		Version: req.GetVersion(),
	})
	if errCr != nil {
		msg := fmt.Sprintf("error creating the text for title %s", req.GetTitle())
		h.log.Error(msg, err)
		return nil, handleError(err, msg)
	}

	return &pb.CreateTextResponse{
		Status: "New text created",
	}, nil
}

func (h *TextHandler) GetText(ctx context.Context, req *pb.GetTextRequest) (*pb.GetTextResponse, error) {
	h.log.Debug("GetText handler")
	if req.GetTitle() == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}

	userID := utils.GetUserIDFromCTX(ctx)

	res, err := h.storage.Texts().GetByTitle(ctx, &dto.GetTextByTitleRequest{
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

	return &pb.GetTextResponse{
		Title:   req.GetTitle(),
		Data:    res.Data,
		Meta:    res.Meta,
		Version: res.Version,
	}, nil
}

func (h *TextHandler) DeleteText(ctx context.Context, req *pb.DeleteTextRequest) (*pb.DeleteTextResponse, error) {
	h.log.Debug("DeleteText handler")
	if req.GetTitle() == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}

	userID := utils.GetUserIDFromCTX(ctx)

	err := h.storage.Texts().Delete(ctx, &dto.DeleteTextRequest{
		Title:  req.GetTitle(),
		UserID: *userID,
	})
	if err != nil {
		msg := fmt.Sprintf("failed to delete data for title %s", req.GetTitle())
		h.log.Error(msg, err)
		return nil, status.Error(codes.Internal, msg)
	}

	return &pb.DeleteTextResponse{
		Status: "success",
	}, nil

}
