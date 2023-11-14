package worker

import (
	"context"

	"github.com/mauriciofsnts/hermes/internal/types"
)

var cancel context.CancelFunc

func StartWorker(queue types.Queue[types.Email]) {
	var ctx context.Context

	ctx, cancel = context.WithCancel(context.Background())
	go queue.Read(ctx)
}

func StopWorker() {
	cancel()
}
