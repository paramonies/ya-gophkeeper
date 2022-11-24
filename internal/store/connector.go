package store

import (
	"context"

	"github.com/paramonies/ya-gophkeeper/internal/store/dto"
)

type Connector interface {
	Users() UserRepo
	Passwords() PasswordRepo
}

type UserRepo interface {
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error)
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error)
}

type PasswordRepo interface {
	Create(ctx context.Context, req *dto.CreateRequest) (*dto.CreateResponse, error)
	GetByLogin(ctx context.Context, req *dto.GetByLoginRequest) (*dto.GetByLoginResponse, error)
	GetAll(ctx context.Context, req *dto.GetAllRequest) (*dto.GetAllResponse, error)
	Delete(ctx context.Context, req *dto.DeletePasswordRequest) error
}
