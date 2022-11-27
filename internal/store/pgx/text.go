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

type TextRepo struct {
	pool         *pgxpool.Pool
	queryTimeout time.Duration
}

func NewTextRepo(p *pgxpool.Pool, queryTimeout time.Duration) *TextRepo {
	return &TextRepo{
		pool:         p,
		queryTimeout: queryTimeout,
	}
}

func (r *TextRepo) Create(ctx context.Context, req *dto.CreateTextRequest) (*dto.CreateTextResponse, error) {
	query := `
INSERT INTO texts
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

	var textID string
	row := r.pool.QueryRow(ctx, query, req.UserID, req.Title, req.Data, req.Meta, req.Version)

	if err := row.Scan(&textID); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == pgerrcode.UniqueViolation {
			return nil, core.NewUniqueViolationError(pgErr.ConstraintName, fmt.Sprintf("text with title %s already exists", req.Title))
		}

		return nil, err
	}

	return &dto.CreateTextResponse{
		TextID: textID,
	}, nil
}

func (r *TextRepo) GetByLogin(ctx context.Context, req *dto.GetTextByTitleRequest) (*dto.GetTextByTitleResponse, error) {
	query := `
SELECT title, data, meta, version FROM texts
WHERE title=$1 AND user_id=$2 AND deleted_at isnull ORDER BY version DESC LIMIT 1
`

	ctx, cancel := context.WithTimeout(ctx, r.queryTimeout)
	defer cancel()

	var (
		id      string
		data    string
		meta    string
		version uint32
	)
	err := r.pool.QueryRow(ctx, query, req.Title, req.UserID).Scan(&id, &data, &meta, &version)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, core.NewTextNotFoundError(req.Title)
		}
		return nil, err
	}

	return &dto.GetTextByTitleResponse{
		ID:      id,
		UserID:  req.UserID,
		Title:   req.Title,
		Data:    data,
		Meta:    meta,
		Version: version,
	}, nil
}

func (r *TextRepo) GetAll(ctx context.Context, req *dto.GetTextAllRequest) (*dto.GetTextAllResponse, error) {
	query := `
SELECT DISTINCT ON (title) title, data, meta, version 
FROM texts WHERE user_id = $1 AND deleted_at isnull ORDER BY title, version DESC
`

	ctx, cancel := context.WithTimeout(ctx, r.queryTimeout)
	defer cancel()

	rows, err := r.pool.Query(ctx, query, req.UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	texts := make([]*model.Text, 0)
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

		t := &model.Text{
			Title:   title,
			Data:    data,
			Meta:    meta,
			Version: version,
		}
		texts = append(texts, t)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return &dto.GetTextAllResponse{
		Texts: texts,
	}, nil
}

func (r *TextRepo) Delete(ctx context.Context, req *dto.DeleteTextRequest) error {
	query := `
UPDATE texts SET deleted_at = current_timestamp WHERE title = $1 AND user_id = $2
`
	ctx, cancel := context.WithTimeout(ctx, r.queryTimeout)
	defer cancel()

	if _, err := r.pool.Exec(ctx, query, req.Title, req.UserID); err != nil {
		return fmt.Errorf("error deleting text for %s title: %w", req.Title, err)
	}

	return nil
}
