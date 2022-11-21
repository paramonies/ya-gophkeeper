package store

import (
	"context"

	"github.com/paramonies/ya-gophkeeper/internal/store/dto"
)

type Connector interface {
	Users() UserRepo
	//Passwords() PasswordRepo
}

type UserRepo interface {
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error)
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error)
}

type PasswordRepo interface {
	Create(ctx context.Context, req *dto.CreatePasswordRequest) (*dto.CreatePasswordResponse, error)
	Get(ctx context.Context, req *dto.GetPasswordRequest) (*dto.GetPasswordResponse, error)
	Update(ctx context.Context, req *dto.UpdatePasswordRequest) (*dto.UpdatePasswordResponse, error)
	Delete(ctx context.Context, req *dto.DeletePasswordRequest) (*dto.DeletePasswordResponse, error)
}
