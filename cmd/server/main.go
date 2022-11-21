package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/paramonies/ya-gophkeeper/internal/server"
	"github.com/paramonies/ya-gophkeeper/internal/server/config"
	"github.com/paramonies/ya-gophkeeper/internal/store"
	"github.com/paramonies/ya-gophkeeper/pkg/logger"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("failed to parse client config", err)
	}

	l := logger.New(cfg.Log.Level)
	l.Info("start gophkeeper server")
	dbPool, err := initDatabaseConnection(cfg.DB.DNS, time.Duration(cfg.DB.ConnectTimeout)*time.Second, l)
	dbConn := store.NewPgxConnector(dbPool, time.Duration(cfg.DB.QueryTimeout)*time.Second)

	// init the grpc server
	server, err := server.InitGRPCServer(dbConn, l)
	if err != nil {
		log.Fatal(err)
	}

	// init storage
	//if err = storage.Init(); err != nil {
	//	log.Fatal(err)
	//}

	gracefulShutdownChan := make(chan struct{})
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		<-sigint
		l.Info("server gracefully shutdown: start")
		if err = server.ShutDown(); err != nil {
			l.Info("gRPC server shutdown err: %v", err)
		}
		close(gracefulShutdownChan)
	}()

	if err = server.Start(cfg.Server.Address); err != nil {
		log.Fatal(err)
	}

	<-gracefulShutdownChan

	//if err = storage.Close(); err != nil {
	//	log.Fatal(err)
	//}

	l.Info("server gracefully shutdown: done")
}

func initDatabaseConnection(dns string, dbConnectTimeout time.Duration, l *logger.Logger) (*pgxpool.Pool, error) {
	l.Info("init dbase")
	ctx, cancel := context.WithTimeout(context.Background(), dbConnectTimeout)
	defer cancel()

	pool, err := pgxpool.Connect(ctx, dns)
	if err != nil {
		return nil, err
	}

	return pool, nil
}
