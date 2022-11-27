package service

import (
	"context"
	"fmt"
	"unicode/utf8"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/paramonies/ya-gophkeeper/internal/core"
	"github.com/paramonies/ya-gophkeeper/internal/server/utils"
	"github.com/paramonies/ya-gophkeeper/internal/store"
	"github.com/paramonies/ya-gophkeeper/internal/store/dto"
	pb "github.com/paramonies/ya-gophkeeper/pkg/gen/api/gophkeeper/v1"
	"github.com/paramonies/ya-gophkeeper/pkg/logger"
)

type CardHandler struct {
	pb.UnimplementedCardServiceServer
	storage store.Connector
	log     *logger.Logger
}

func NewCardHandler(s store.Connector, l *logger.Logger) *CardHandler {
	return &CardHandler{
		storage: s,
		log:     l,
	}
}

func (h *CardHandler) CreateCard(ctx context.Context, req *pb.CreateCardRequest) (*pb.CreateCardResponse, error) {
	h.log.Debug("CreateCard handler")
	cvvLen := utf8.RuneCountInString(string(req.GetCvv()))
	if req.GetVersion() < 1 || req.GetNumber() == "" || req.GetOwner() == "" || req.GetExpDate() == "" || cvvLen != 3 {
		return nil, status.Error(codes.InvalidArgument, "invalid arguments")
	}

	userID := utils.GetUserIDFromCTX(ctx)

	res, err := h.storage.Cards().GetByNumber(ctx, &dto.GetCardByNumberRequest{
		Number: req.GetNumber(),
		UserID: *userID,
	})
	if err != nil && !core.IsNotFound(err) {
		msg := fmt.Sprintf("error getting the card with %s number from db", req.GetNumber())
		h.log.Error(msg, err)
		return nil, handleError(err, msg)
	}

	if err != nil && core.IsNotFound(err) {
		_, errCr := h.storage.Cards().Create(ctx, &dto.CreateCardRequest{
			UserID:  *userID,
			Number:  req.GetNumber(),
			Owner:   req.GetOwner(),
			ExpDate: req.GetExpDate(),
			Cvv:     req.GetCvv(),
			Meta:    req.GetMeta(),
			Version: req.GetVersion(),
		})
		if errCr != nil {
			msg := fmt.Sprintf("error creating the card with number %s", req.GetNumber())
			h.log.Error(msg, errCr)
			return nil, handleError(errCr, msg)
		}

		return &pb.CreateCardResponse{
			Status: "New card created",
		}, nil
	}

	if req.GetVersion() <= res.Version {
		return nil, status.Error(codes.AlreadyExists, newerVersionDetected)
	}

	_, errCr := h.storage.Cards().Create(ctx, &dto.CreateCardRequest{
		UserID:  *userID,
		Number:  req.GetNumber(),
		Owner:   req.GetOwner(),
		ExpDate: req.GetExpDate(),
		Cvv:     req.GetCvv(),
		Meta:    req.GetMeta(),
		Version: req.GetVersion(),
	})
	if errCr != nil {
		msg := fmt.Sprintf("error creating the card with number %s", req.GetNumber())
		h.log.Error(msg, err)
		return nil, handleError(err, msg)
	}

	return &pb.CreateCardResponse{
		Status: "New card created",
	}, nil
}

func (h *CardHandler) GetCard(ctx context.Context, req *pb.GetCardRequest) (*pb.GetCardResponse, error) {
	h.log.Debug("GetCard handler")
	if req.GetNumber() == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}

	userID := utils.GetUserIDFromCTX(ctx)

	res, err := h.storage.Cards().GetByNumber(ctx, &dto.GetCardByNumberRequest{
		Number: req.GetNumber(),
		UserID: *userID,
	})
	if err != nil {
		if core.IsNotFound(err) {
			msg := fmt.Sprintf("data for card %s not found", req.GetNumber())
			h.log.Error(msg, err)
			return nil, handleError(err, msg)
		}

		msg := fmt.Sprintf("failed to obtain latest data for card %s", req.GetNumber())
		h.log.Error(msg, err)
		return nil, handleError(err, msg)
	}

	return &pb.GetCardResponse{
		Number:  req.GetNumber(),
		Owner:   res.Owner,
		ExpDate: res.ExpDate,
		Cvv:     res.Cvv,
		Meta:    res.Meta,
		Version: res.Version,
	}, nil
}

func (h *CardHandler) DeleteCard(ctx context.Context, req *pb.DeleteCardRequest) (*pb.DeleteCardResponse, error) {
	h.log.Debug("DeleteCard handler")
	if req.GetNumber() == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}

	userID := utils.GetUserIDFromCTX(ctx)

	err := h.storage.Cards().Delete(ctx, &dto.DeleteCardRequest{
		Number: req.GetNumber(),
		UserID: *userID,
	})
	if err != nil {
		msg := fmt.Sprintf("failed to delete data for card %s", req.GetNumber())
		h.log.Error(msg, err)
		return nil, status.Error(codes.Internal, msg)
	}

	return &pb.DeleteCardResponse{
		Status: "success",
	}, nil

}
