package dto

import "github.com/paramonies/ya-gophkeeper/internal/model"

type CreateTextRequest struct {
	UserID  string
	Title   string
	Data    string
	Meta    string
	Version uint32
}

type CreateTextResponse struct {
	TextID string
}

type GetTextByTitleRequest struct {
	Title  string
	UserID string
}

type GetTextByTitleResponse struct {
	ID      string
	UserID  string
	Title   string
	Data    string
	Meta    string
	Version uint32
}

type GetTextAllRequest struct {
	UserID string
}

type GetTextAllResponse struct {
	Texts []*model.Text
}
type DeleteTextRequest struct {
	Title  string
	UserID string
}
