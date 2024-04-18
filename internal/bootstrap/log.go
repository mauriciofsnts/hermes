package bootstrap

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"github.com/mauriciofsnts/hermes/internal/config"
)

func setupLog(cfg *config.Config) {
	var handler slog.Handler
	level := cfg.Log.Level
	showSource := cfg.Log.ShowSource

	switch cfg.Log.Type {
	case config.LogTypeJSON:
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     level,
			AddSource: showSource,
		})
	case config.LogTypeColored:
		handler = tint.NewHandler(os.Stdout, &tint.Options{
			Level:      level,
			TimeFormat: time.DateTime,
			AddSource:  showSource,
		})
	default:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level:     level,
			AddSource: showSource,
		})
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
}
