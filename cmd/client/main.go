package main

import (
	"github.com/paramonies/ya-gophkeeper/pkg/logger"
)

func main() {
	l := logger.New("debug")
	l.Info("hello")
}
