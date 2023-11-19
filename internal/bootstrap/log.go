package bootstrap

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
)

func SetupLog() {
	handler := tint.NewHandler(os.Stdout, &tint.Options{
		TimeFormat: time.DateTime,
	})

	logger := slog.New(handler)

	slog.SetDefault(logger)
}
