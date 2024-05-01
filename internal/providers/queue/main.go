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
var Cancel context.CancelFunc

func NewQueue(cfg *config.Config) error {
	if cfg.Redis != nil {
		redisQueue, err := redis.NewRedisProvider()

		if err == nil {
			Queue = redisQueue
			return nil
		}
	}

	slog.Warn("Using memory queue, because no queue provider was found")
	memoryQueue := memory.NewMemoryProvider()
	Queue = memoryQueue

	return nil
}

func StartWorker() {
	var ctx context.Context

	ctx, Cancel = context.WithCancel(context.Background())
	go Queue.Read(ctx)
}

func StopWorker() {
	Cancel()
}
