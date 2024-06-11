package queue

import (
	"context"
	"log/slog"

	"github.com/mauriciofsnts/hermes/internal/config"
	"github.com/mauriciofsnts/hermes/internal/providers/queue/memory"
	"github.com/mauriciofsnts/hermes/internal/providers/queue/redis"
	"github.com/mauriciofsnts/hermes/internal/providers/queue/worker"
	"github.com/mauriciofsnts/hermes/internal/types"
)

var Cancel context.CancelFunc

func NewQueue(cfg *config.Config) (worker.Queue[types.Mail], error) {
	if cfg.Redis != nil {
		redisQueue, err := redis.NewRedisProvider()

		if err == nil {
			return redisQueue, nil
		}
	}

	slog.Warn("Using memory queue, because no queue provider was found")
	memoryQueue := memory.NewMemoryProvider()
	return memoryQueue, nil
}

func StartWorker(queue worker.Queue[types.Mail]) {
	var ctx context.Context

	ctx, Cancel = context.WithCancel(context.Background())
	go queue.Read(ctx)
}

func StopWorker() {
	Cancel()
}
