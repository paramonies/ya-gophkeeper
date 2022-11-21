package pgx

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/paramonies/ya-gophkeeper/internal/core"
	"github.com/paramonies/ya-gophkeeper/internal/store/dto"
)

type UserRepo struct {
	pool         *pgxpool.Pool
	queryTimeout time.Duration
}

func NewUserRepo(p *pgxpool.Pool, queryTimeout time.Duration) *UserRepo {
	return &UserRepo{
		pool:         p,
		queryTimeout: queryTimeout,
	}
}

func (r *UserRepo) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	query := `
INSERT INTO users
(
 login,
 password_hash
)
VALUES ($1, $2)
RETURNING id
`
	ctx, cancel := context.WithTimeout(ctx, r.queryTimeout)
	defer cancel()

	var userID string
	row := r.pool.QueryRow(ctx, query, req.Login, req.PasswordHash)

	if err := row.Scan(&userID); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == pgerrcode.UniqueViolation {
			return nil, core.NewUniqueViolationError(pgErr.ConstraintName, fmt.Sprintf("such login %s already exists", req.Login))
		}

		return nil, err
	}

	return &dto.RegisterResponse{
		UserID: userID,
	}, nil
}
