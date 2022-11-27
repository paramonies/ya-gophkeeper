package store

import (
	"context"

	"github.com/paramonies/ya-gophkeeper/internal/store/dto"
)

type Connector interface {
	Users() UserRepo
	Passwords() PasswordRepo
	Texts() TextRepo
	Binaries() BinaryRepo
}

type UserRepo interface {
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error)
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error)
}

type PasswordRepo interface {
	Create(ctx context.Context, req *dto.CreatePwdRequest) (*dto.CreatePwdResponse, error)
	GetByLogin(ctx context.Context, req *dto.GetPwdByLoginRequest) (*dto.GetPwdByLoginResponse, error)
	GetAll(ctx context.Context, req *dto.GetPwdAllRequest) (*dto.GetPwdAllResponse, error)
	Delete(ctx context.Context, req *dto.DeletePwdRequest) error
}

type TextRepo interface {
	Create(ctx context.Context, req *dto.CreateTextRequest) (*dto.CreateTextResponse, error)
	GetByLogin(ctx context.Context, req *dto.GetTextByTitleRequest) (*dto.GetTextByTitleResponse, error)
	GetAll(ctx context.Context, req *dto.GetTextAllRequest) (*dto.GetTextAllResponse, error)
	Delete(ctx context.Context, req *dto.DeleteTextRequest) error
}

type BinaryRepo interface {
	Create(ctx context.Context, req *dto.CreateBinaryRequest) (*dto.CreateBinaryResponse, error)
	GetByLogin(ctx context.Context, req *dto.GetBinaryByTitleRequest) (*dto.GetBinaryByTitleResponse, error)
	GetAll(ctx context.Context, req *dto.GetBinaryAllRequest) (*dto.GetBinaryAllResponse, error)
	Delete(ctx context.Context, req *dto.DeleteBinaryRequest) error
}
