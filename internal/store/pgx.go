package store

import (
	"time"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/paramonies/ya-gophkeeper/internal/store/pgx"
)

type pgxConnector struct {
	userRepo     *pgx.UserRepo
	passwordRepo *pgx.PasswordRepo
	textRepo     *pgx.TextRepo
	binaryRepo   *pgx.BinaryRepo
}

func NewPgxConnector(p *pgxpool.Pool, queryTimeout time.Duration) Connector {
	return &pgxConnector{
		userRepo:     pgx.NewUserRepo(p, queryTimeout),
		passwordRepo: pgx.NewPasswordRepo(p, queryTimeout),
		textRepo:     pgx.NewTextRepo(p, queryTimeout),
		binaryRepo:   pgx.NewBinaryRepo(p, queryTimeout),
	}
}

func (c *pgxConnector) Users() UserRepo {
	return c.userRepo
}

func (c *pgxConnector) Passwords() PasswordRepo {
	return c.passwordRepo
}

func (c *pgxConnector) Texts() TextRepo {
	return c.textRepo
}

func (c *pgxConnector) Binaries() BinaryRepo {
	return c.binaryRepo
}
