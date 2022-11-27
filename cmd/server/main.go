package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/paramonies/ya-gophkeeper/pkg/graceful"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/paramonies/ya-gophkeeper/internal/server"
	"github.com/paramonies/ya-gophkeeper/internal/server/config"
	"github.com/paramonies/ya-gophkeeper/internal/store"
	"github.com/paramonies/ya-gophkeeper/pkg/logger"
)

const errorExitCode int = 1

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("failed to parse client config", err)
	}

	l := logger.New(cfg.Log.Level)
	l.Info("start gophkeeper server")
	dbPool, err := initDatabaseConnection(cfg.DB.DNS, time.Duration(cfg.DB.ConnectTimeout)*time.Second, l)
	dbConn := store.NewPgxConnector(dbPool, time.Duration(cfg.DB.QueryTimeout)*time.Second)

	err = server.RunGRPCServer(cfg.Server.Address, dbConn, l)
	if err != nil {
		l.Error(context.Background(), "failed to run API server", err)
		os.Exit(errorExitCode)
	}

	graceful.AddCallback(func() error {
		l.Warn("Shutting down application...")
		return nil
	})

	err = graceful.WaitShutdown()
	if err != nil {
		l.Error(context.Background(), "shutdown error", err)
	}
}

func initDatabaseConnection(dns string, dbConnectTimeout time.Duration, l *logger.Logger) (*pgxpool.Pool, error) {
	l.Info("init dbase")
	ctx, cancel := context.WithTimeout(context.Background(), dbConnectTimeout)
	defer cancel()

	pool, err := pgxpool.Connect(ctx, dns)
	if err != nil {
		return nil, err
	}

	graceful.AddCallback(func() error {
		pool.Close()
		return nil
	})

	return pool, nil
}
