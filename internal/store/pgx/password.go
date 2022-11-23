package pgx

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/paramonies/ya-gophkeeper/internal/core"
	"github.com/paramonies/ya-gophkeeper/internal/store/dto"
)

type PasswordRepo struct {
	pool         *pgxpool.Pool
	queryTimeout time.Duration
}

func NewPasswordRepo(p *pgxpool.Pool, queryTimeout time.Duration) *PasswordRepo {
	return &PasswordRepo{
		pool:         p,
		queryTimeout: queryTimeout,
	}
}

func (r *PasswordRepo) Create(ctx context.Context, req *dto.CreateRequest) (*dto.CreateResponse, error) {
	query := `
INSERT INTO passwords
(
 user_id,
 login,
 password,
 meta,
 version
)
VALUES ($1, $2, $3, $4, $5)
RETURNING id
`
	ctx, cancel := context.WithTimeout(ctx, r.queryTimeout)
	defer cancel()

	var passID string
	row := r.pool.QueryRow(ctx, query, req.UserID, req.Login, req.Password, req.Meta, req.Version)

	if err := row.Scan(&passID); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == pgerrcode.UniqueViolation {
			return nil, core.NewUniqueViolationError(pgErr.ConstraintName, fmt.Sprintf("password with login %s already exists", req.Login))
		}

		return nil, err
	}

	return &dto.CreateResponse{
		PasswordID: passID,
	}, nil
}

func (r *PasswordRepo) GetByID(ctx context.Context, req *dto.GetByIDRequest) (*dto.GetByIDResponse, error) {
	query := `
SELECT login, password, meta, version FROM passwords
WHERE login=$1 AND user_id=$2 AND deleted_at isnull ORDER BY version DESC LIMIT 1
`

	ctx, cancel := context.WithTimeout(ctx, r.queryTimeout)
	defer cancel()

	var (
		id       string
		password string
		meta     string
		version  uint32
	)
	err := r.pool.QueryRow(ctx, query, req.Login, req.UserID).Scan(&id, &password, &meta, &version)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, core.NewPasswordNotFoundError(req.Login)
		}
		return nil, err
	}

	return &dto.GetByIDResponse{
		ID:       id,
		UserID:   req.UserID,
		Login:    req.Login,
		Password: password,
		Meta:     meta,
		Version:  version,
	}, nil
}
