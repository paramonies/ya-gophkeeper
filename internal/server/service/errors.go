package service

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/paramonies/ya-gophkeeper/internal/core"
)

func handleError(err error, msg string) error {
	if err == nil {
		return nil
	}

	switch {
	case core.IsUniqueViolationError(err):
		return status.Error(codes.InvalidArgument, err.Error())
	case core.IsNotFound(err):
		return status.Error(codes.NotFound, err.Error())
	}

	return status.Error(codes.Internal, msg)
}
