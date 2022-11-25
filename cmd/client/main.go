package main

import (
	"log"
	"os"

	"github.com/paramonies/ya-gophkeeper/internal/client"
	"github.com/paramonies/ya-gophkeeper/internal/client/cmd"
	"github.com/paramonies/ya-gophkeeper/internal/client/config"
	"github.com/paramonies/ya-gophkeeper/internal/client/storage"
	"github.com/paramonies/ya-gophkeeper/pkg/logger"
)

const errorExitCode int = 1

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("failed to parse client config", err)
	}

	l := logger.New(cfg.Log.Level)
	l.Info("start gophkeeper client")

	err = storage.InitStorage(cfg.UsersStoragePath, cfg.ObjectsStoragePath)
	if err != nil {
		l.Error("failed to init client local storage ", err)
		os.Exit(errorExitCode)
	}
	l.Info("local storage initialized")

	clientSet, err := client.CreateClientSet(cfg.Server.GrpcServerPath)
	if err != nil {
		l.Error("failed to initiate a connection to service", err)
		os.Exit(errorExitCode)
	}
	defer client.ConnDown()

	err = cmd.Init(l, cfg, clientSet)
	if err != nil {
		l.Error("failed to init cobra commands", err)
		os.Exit(errorExitCode)
	}

}
