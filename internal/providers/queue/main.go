package queue

import (
	"context"
	"log/slog"

	"github.com/mauriciofsnts/hermes/internal/config"
	"github.com/mauriciofsnts/hermes/internal/providers/queue/memory"
	"github.com/mauriciofsnts/hermes/internal/providers/queue/redis"
	"github.com/mauriciofsnts/hermes/internal/types"
)

var Queue types.Queue[types.Mail]
var cancel context.CancelFunc

func NewQueue(cfg *config.Config) (types.Queue[types.Mail], error) {
	if cfg.Redis != nil {
		Queue, err := redis.NewRedisProvider()

		if err == nil {
			return Queue, nil
		}
	}

	slog.Warn("Using memory queue, because no queue provider was found")
	memoryQueue := memory.NewMemoryProvider()

	return memoryQueue, nil
}

func StartWorker(queue types.Queue[types.Mail]) {
	var ctx context.Context

	ctx, cancel = context.WithCancel(context.Background())
	go queue.Read(ctx)
}

func StopWorker() {
	cancel()
}
