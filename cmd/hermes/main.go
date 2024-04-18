package main

import (
	"log/slog"
	"os"

	"github.com/mauriciofsnts/hermes/internal/bootstrap"
	"github.com/mauriciofsnts/hermes/internal/config"
)

const (
	DefaultConfigPath = "config.yaml"
)

func main() {
	cfg, err := config.LoadConfigFromFile(DefaultConfigPath)

	if err != nil {
		slog.Error("failed to load config file: %v", err)
		os.Exit(1)

	}

	bootstrap.Start(cfg)
}
