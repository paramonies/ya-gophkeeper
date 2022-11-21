package pgx

import (
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
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
