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

type CardRepo struct {
	pool         *pgxpool.Pool
	queryTimeout time.Duration
}

func NewCardRepo(p *pgxpool.Pool, queryTimeout time.Duration) *CardRepo {
	return &CardRepo{
		pool:         p,
		queryTimeout: queryTimeout,
	}
}

func (r *CardRepo) Create(ctx context.Context, req *dto.CreateCardRequest) (*dto.CreateCardResponse, error) {
	query := `
INSERT INTO cards
(
 user_id,
 number,
 owner,
 expiration_date,
 cvv,
 meta,
 version
)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id
`
	ctx, cancel := context.WithTimeout(ctx, r.queryTimeout)
	defer cancel()

	var cardID string
	row := r.pool.QueryRow(ctx, query, req.UserID, req.Number, req.Owner, req.ExpDate, req.Cvv, req.Meta, req.Version)

	if err := row.Scan(&cardID); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == pgerrcode.UniqueViolation {
			return nil, core.NewUniqueViolationError(pgErr.ConstraintName, fmt.Sprintf("card with number %s already exists", req.Number))
		}

		return nil, err
	}

	return &dto.CreateCardResponse{
		CardID: cardID,
	}, nil
}

func (r *CardRepo) GetByNumber(ctx context.Context, req *dto.GetCardByNumberRequest) (*dto.GetCardByNumberResponse, error) {
	query := `
SELECT number, owner, expiration_date, cvv, meta, version FROM cards
WHERE number=$1 AND user_id=$2 AND deleted_at isnull ORDER BY version DESC LIMIT 1
`

	ctx, cancel := context.WithTimeout(ctx, r.queryTimeout)
	defer cancel()

	var (
		id      string
		owner   string
		expDate string
		cvv     string
		meta    string
		version uint32
	)
	err := r.pool.QueryRow(ctx, query, req.Number, req.UserID).Scan(&id, &owner, &expDate, &cvv, &meta, &version)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, core.NewCardNotFoundError(req.Number)
		}
		return nil, err
	}

	return &dto.GetCardByNumberResponse{
		ID:      id,
		UserID:  req.UserID,
		Number:  req.Number,
		Owner:   owner,
		ExpDate: expDate,
		Cvv:     cvv,
		Meta:    meta,
		Version: version,
	}, nil
}

func (r *CardRepo) GetAll(ctx context.Context, req *dto.GetCardAllRequest) (*dto.GetCardAllResponse, error) {
	query := `
SELECT DISTINCT ON (number) number, owner, expiration_date, cvv, meta, version 
FROM cards WHERE user_id = $1 AND deleted_at isnull ORDER BY number, version DESC
`

	ctx, cancel := context.WithTimeout(ctx, r.queryTimeout)
	defer cancel()

	rows, err := r.pool.Query(ctx, query, req.UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cards := make([]*model.Card, 0)
	for rows.Next() {
		var (
			number  string
			owner   string
			expDate string
			cvv     string
			meta    string
			version uint32
		)

		err = rows.Scan(&number, &owner, &expDate, &cvv, &meta, &version)
		if err != nil {
			return nil, err
		}

		c := &model.Card{
			Number:  number,
			Owner:   owner,
			ExpDate: expDate,
			Cvv:     cvv,
			Meta:    meta,
			Version: version,
		}
		cards = append(cards, c)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return &dto.GetCardAllResponse{
		Cards: cards,
	}, nil
}

func (r *CardRepo) Delete(ctx context.Context, req *dto.DeleteCardRequest) error {
	query := `
UPDATE cards SET deleted_at = current_timestamp WHERE number = $1 AND user_id = $2
`
	ctx, cancel := context.WithTimeout(ctx, r.queryTimeout)
	defer cancel()

	if _, err := r.pool.Exec(ctx, query, req.Number, req.UserID); err != nil {
		return fmt.Errorf("error deleting card with %s number: %w", req.Number, err)
	}

	return nil
}
