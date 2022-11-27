package dto

import "github.com/paramonies/ya-gophkeeper/internal/model"

type CreateBinaryRequest struct {
	UserID  string
	Title   string
	Data    string
	Meta    string
	Version uint32
}

type CreateBinaryResponse struct {
	BianryID string
}

type GetBinaryByTitleRequest struct {
	Title  string
	UserID string
}

type GetBinaryByTitleResponse struct {
	ID      string
	UserID  string
	Title   string
	Data    string
	Meta    string
	Version uint32
}

type GetBinaryAllRequest struct {
	UserID string
}

type GetBinaryAllResponse struct {
	Binaries []*model.Binary
}
type DeleteBinaryRequest struct {
	Title  string
	UserID string
}
