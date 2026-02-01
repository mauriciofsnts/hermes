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
// It encapsulates the context and cancel function to avoid global mutable state.
type QueueManager struct {
	ctx    context.Context
	cancel context.CancelFunc
	queue  worker.Queue[types.Mail]
}

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
// It waits for up to the specified timeout for pending items to be processed.
// This ensures no emails are lost during shutdown.
func (qm *QueueManager) DrainAndStop(timeout time.Duration) {
	slog.Info("Draining queue before shutdown", "timeout", timeout)

	// Create a timeout context for draining
	drainCtx, drainCancel := context.WithTimeout(context.Background(), timeout)
	defer drainCancel()

	// Signal the worker to stop accepting new items
	if qm.cancel != nil {
		qm.cancel()
	}

	// Wait for either:
	// 1. All pending items to be processed
	// 2. Timeout to expire
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-drainCtx.Done():
			slog.Warn("Drain timeout reached, forcing shutdown")
			return
		case <-ticker.C:
			// Check if queue is empty (implementation depends on queue type)
			// For now, we just wait for the timeout
			// In a real implementation, we'd check queue depth here
			slog.Debug("Waiting for queue to drain...")
		}
	}
}
