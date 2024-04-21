package queue

import (
	"context"

	"github.com/mauriciofsnts/hermes/internal/config"
	"github.com/mauriciofsnts/hermes/internal/providers/queue/kafka"
	"github.com/mauriciofsnts/hermes/internal/providers/queue/memory"
	"github.com/mauriciofsnts/hermes/internal/providers/queue/redis"
	"github.com/mauriciofsnts/hermes/internal/types"
)

var Queue types.Queue[types.Mail]
var cancel context.CancelFunc

func NewQueue(cfg *config.Config) (types.Queue[types.Mail], error) {
	if cfg.Kafka != nil {
		err := kafka.CreateTopic()

		if err != nil {
			return nil, err
		}

		Queue = kafka.NewKafkaProvider()
	} else if cfg.Redis != nil {
		Queue = redis.NewRedisProvider()
	} else {
		Queue = memory.NewMemoryProvider()
	}

	return Queue, nil
}

func StartWorker(queue types.Queue[types.Mail]) {
	var ctx context.Context

	ctx, cancel = context.WithCancel(context.Background())
	go queue.Read(ctx)
}

func StopWorker() {
	cancel()
}
