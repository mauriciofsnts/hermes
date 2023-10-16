package worker

import (
	"context"

	"github.com/mauriciofsnts/hermes/internal/storage"
)

var cancel context.CancelFunc

func StartWorker() {
	storage := storage.NewStorage()

	var ctx context.Context

	ctx, cancel = context.WithCancel(context.Background())

	go storage.Read(ctx)
}

func StopWorker() {
	cancel()
}
