package dto

import "github.com/paramonies/ya-gophkeeper/internal/model"

type CreateRequest struct {
	UserID   string
	Login    string
	Password string
	Meta     string
	Version  uint32
}

type CreateResponse struct {
	PasswordID string
}

type GetByLoginRequest struct {
	Login  string
	UserID string
}

type GetByLoginResponse struct {
	ID       string
	UserID   string
	Login    string
	Password string
	Meta     string
	Version  uint32
}

type GetAllRequest struct {
	UserID string
}

type GetAllResponse struct {
	Passwords []*model.Password
}
type DeletePasswordRequest struct {
	Login  string
	UserID string
}
