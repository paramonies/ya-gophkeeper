package main

import (
	"log"
	"os"

	"github.com/paramonies/ya-gophkeeper/internal/client/cmd"
	"github.com/paramonies/ya-gophkeeper/internal/client/config"
	client "github.com/paramonies/ya-gophkeeper/internal/client/grpc"
	"github.com/paramonies/ya-gophkeeper/internal/client/storage"
	"github.com/paramonies/ya-gophkeeper/pkg/logger"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("failed to parse client config", err)
	}

	_, err = os.OpenFile(cfg.UsersStoragePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		log.Fatal("failed to create file ", err)
	}

	l := logger.New(cfg.Log.Level)
	l.Info("start gophkeeper client service")

	err = storage.InitStorage(cfg.UsersStoragePath, cfg.ObjectsStoragePath)
	if err != nil {
		log.Fatal("failed to init client local storage ", err)
	}
	l.Info("local storage initialized")

	cliUser, err := client.DialUpUser(cfg.Server.GrpcServerPath)
	if err != nil {
		l.Fatal("failed to initiates a connection to user server", err)
	}

	cliPass, err := client.DialUpPass(cfg.Server.GrpcServerPath)
	if err != nil {
		l.Fatal("failed to initiates a connection to user server", err)
	}

	err = cmd.Init(l, cliUser, cliPass, cfg)
	if err != nil {
		l.Fatal("failed to init cobra commands")
	}

}
