package dto

import "github.com/paramonies/ya-gophkeeper/internal/model"

type CreateCardRequest struct {
	UserID  string
	Number  string
	Owner   string
	ExpDate string
	Cvv     string
	Meta    string
	Version uint32
}

type CreateCardResponse struct {
	CardID string
}

type GetCardByNumberRequest struct {
	Number string
	UserID string
}

type GetCardByNumberResponse struct {
	ID      string
	UserID  string
	Number  string
	Owner   string
	ExpDate string
	Cvv     string
	Meta    string
	Version uint32
}

type GetCardAllRequest struct {
	UserID string
}

type GetCardAllResponse struct {
	Cards []*model.Card
}
type DeleteCardRequest struct {
	Number string
	UserID string
}
