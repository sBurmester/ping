package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/sburmester/ping/pkg/http"
)

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	slog.Info("starting ping server")
	err := http.StartServer()
	if err != nil {
		slog.Default().Error("error starting server", err)
	}
	return err
}
