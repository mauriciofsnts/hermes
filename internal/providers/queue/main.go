package queue

import (
	"context"
	"log/slog"
	"time"

	"github.com/mauriciofsnts/hermes/internal/config"
	"github.com/mauriciofsnts/hermes/internal/providers/queue/memory"
	"github.com/mauriciofsnts/hermes/internal/providers/queue/redis"
	"github.com/mauriciofsnts/hermes/internal/providers/queue/worker"
	"github.com/mauriciofsnts/hermes/internal/types"
)

// QueueManager manages the lifecycle of a notification queue worker.
type QueueManager struct {
	ctx    context.Context
	cancel context.CancelFunc
	queue  worker.Queue[types.Mail]
}

func NewQueue(cfg *config.Config) (worker.Queue[types.Mail], error) {
	if cfg.Redis != nil && cfg.Redis.Address != "" {
		redisQueue, err := redis.NewRedisProvider()
		if err == nil {
			return redisQueue, nil
		}
		slog.Warn("failed to create redis queue, falling back to memory queue", "error", err)
	}

	slog.Warn("using memory queue because no redis provider was configured")
	return memory.NewMemoryProvider(), nil
}

// NewQueueManager creates a new QueueManager instance and starts the worker.
func NewQueueManager(queue worker.Queue[types.Mail]) *QueueManager {
	ctx, cancel := context.WithCancel(context.Background())

	qm := &QueueManager{
		ctx:    ctx,
		cancel: cancel,
		queue:  queue,
	}

	go qm.queue.Read(ctx)
	return qm
}

// Stop gracefully stops the queue worker.
func (qm *QueueManager) Stop() {
	if qm.cancel != nil {
		qm.cancel()
	}
}

// DrainAndStop gracefully drains remaining items and stops the queue worker.
func (qm *QueueManager) DrainAndStop(timeout time.Duration) {
	slog.Info("draining queue before shutdown", "timeout", timeout)

	drainCtx, drainCancel := context.WithTimeout(context.Background(), timeout)
	defer drainCancel()

	if qm.cancel != nil {
		qm.cancel()
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-drainCtx.Done():
			slog.Warn("drain timeout reached, forcing shutdown")
			return
		case <-ticker.C:
			slog.Debug("waiting for queue to drain...")
		}
	}
}
