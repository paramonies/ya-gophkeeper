package dto

import "github.com/paramonies/ya-gophkeeper/internal/model"

type CreatePwdRequest struct {
	UserID   string
	Login    string
	Password string
	Meta     string
	Version  uint32
}

type CreatePwdResponse struct {
	PasswordID string
}

type GetPwdByLoginRequest struct {
	Login  string
	UserID string
}

type GetPwdByLoginResponse struct {
	ID       string
	UserID   string
	Login    string
	Password string
	Meta     string
	Version  uint32
}

type GetPwdAllRequest struct {
	UserID string
}

type GetPwdAllResponse struct {
	Passwords []*model.Password
}
type DeletePwdRequest struct {
	Login  string
	UserID string
}
