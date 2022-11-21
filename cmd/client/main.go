package main

import (
	"log"

	"github.com/paramonies/ya-gophkeeper/internal/client/cmd"
	"github.com/paramonies/ya-gophkeeper/internal/client/config"
	client "github.com/paramonies/ya-gophkeeper/internal/client/grpc"
	"github.com/paramonies/ya-gophkeeper/pkg/logger"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("failed to parse client config", err)
	}

	l := logger.New(cfg.Log.Level)
	l.Info("start gophkeeper client service")

	cli, err := client.DialUp(cfg.Server.GrpcServerPath)
	if err != nil {
		l.Fatal("failed to initiates a connection to server", err)
	}

	err = cmd.Init(l, cli, cfg)
	if err != nil {
		l.Fatal("failed to init cobra commands")
	}

}
