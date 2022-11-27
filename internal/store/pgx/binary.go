package pgx

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/paramonies/ya-gophkeeper/internal/model"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/paramonies/ya-gophkeeper/internal/core"
	"github.com/paramonies/ya-gophkeeper/internal/store/dto"
)

type BinaryRepo struct {
	pool         *pgxpool.Pool
	queryTimeout time.Duration
}

func NewBinaryRepo(p *pgxpool.Pool, queryTimeout time.Duration) *BinaryRepo {
	return &BinaryRepo{
		pool:         p,
		queryTimeout: queryTimeout,
	}
}

func (r *BinaryRepo) Create(ctx context.Context, req *dto.CreateBinaryRequest) (*dto.CreateBinaryResponse, error) {
	query := `
INSERT INTO binaries
(
 user_id,
 title,
 data,
 meta,
 version
)
VALUES ($1, $2, $3, $4, $5)
RETURNING id
`
	ctx, cancel := context.WithTimeout(ctx, r.queryTimeout)
	defer cancel()

	var binaryID string
	row := r.pool.QueryRow(ctx, query, req.UserID, req.Title, req.Data, req.Meta, req.Version)

	if err := row.Scan(&binaryID); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == pgerrcode.UniqueViolation {
			return nil, core.NewUniqueViolationError(pgErr.ConstraintName, fmt.Sprintf("binary with title %s already exists", req.Title))
		}

		return nil, err
	}

	return &dto.CreateBinaryResponse{
		BianryID: binaryID,
	}, nil
}

func (r *BinaryRepo) GetByLogin(ctx context.Context, req *dto.GetBinaryByTitleRequest) (*dto.GetBinaryByTitleResponse, error) {
	query := `
SELECT title, data, meta, version FROM binaries
WHERE title=$1 AND user_id=$2 AND deleted_at isnull ORDER BY version DESC LIMIT 1
`

	ctx, cancel := context.WithTimeout(ctx, r.queryTimeout)
	defer cancel()

	var (
		id      string
		data    []byte
		meta    string
		version uint32
	)
	err := r.pool.QueryRow(ctx, query, req.Title, req.UserID).Scan(&id, &data, &meta, &version)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, core.NewBinaryNotFoundError(req.Title)
		}
		return nil, err
	}

	return &dto.GetBinaryByTitleResponse{
		ID:      id,
		UserID:  req.UserID,
		Title:   req.Title,
		Data:    string(data),
		Meta:    meta,
		Version: version,
	}, nil
}

func (r *BinaryRepo) GetAll(ctx context.Context, req *dto.GetBinaryAllRequest) (*dto.GetBinaryAllResponse, error) {
	query := `
SELECT DISTINCT ON (title) title, data, meta, version 
FROM binaries WHERE user_id = $1 AND deleted_at isnull ORDER BY title, version DESC
`

	ctx, cancel := context.WithTimeout(ctx, r.queryTimeout)
	defer cancel()

	rows, err := r.pool.Query(ctx, query, req.UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	banaries := make([]*model.Binary, 0)
	for rows.Next() {
		var (
			title   string
			data    string
			meta    string
			version uint32
		)

		err = rows.Scan(&title, &data, &meta, &version)
		if err != nil {
			return nil, err
		}

		bin := &model.Binary{
			Title:   title,
			Data:    data,
			Meta:    meta,
			Version: version,
		}
		banaries = append(banaries, bin)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return &dto.GetBinaryAllResponse{
		Binaries: banaries,
	}, nil
}

func (r *BinaryRepo) Delete(ctx context.Context, req *dto.DeleteBinaryRequest) error {
	query := `
UPDATE binaries SET deleted_at = current_timestamp WHERE title = $1 AND user_id = $2
`
	ctx, cancel := context.WithTimeout(ctx, r.queryTimeout)
	defer cancel()

	if _, err := r.pool.Exec(ctx, query, req.Title, req.UserID); err != nil {
		return fmt.Errorf("error deleting binary for %s title: %w", req.Title, err)
	}

	return nil
}
