package worker

import (
	"context"

	"github.com/mauriciofsnts/hermes/internal/types"
)

var cancel context.CancelFunc

func StartWorker(storage types.Storage[types.Email]) {

	var ctx context.Context

	ctx, cancel = context.WithCancel(context.Background())

	go storage.Read(ctx)
}

func StopWorker() {
	cancel()
}
