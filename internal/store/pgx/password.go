package pgx

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
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
