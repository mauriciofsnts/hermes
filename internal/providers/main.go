package providers

import (
	"context"

	"github.com/mauriciofsnts/hermes/internal/config"
	"github.com/mauriciofsnts/hermes/internal/providers/kafka"
	"github.com/mauriciofsnts/hermes/internal/providers/memory"
	"github.com/mauriciofsnts/hermes/internal/providers/redis"
	"github.com/mauriciofsnts/hermes/internal/types"
)

var Queue types.Queue[types.Mail]
var cancel context.CancelFunc

func NewQueue() types.Queue[types.Mail] {
	kafkaEnabled := config.Envs.Kafka.Enabled
	redisEnabled := config.Envs.Redis.Enabled

	if kafkaEnabled {
		err := kafka.CreateTopic()

		if err != nil {
			panic(err)
		}

		Queue = kafka.NewKafkaProvider()
	} else if redisEnabled {
		Queue = redis.NewRedisProvider()
	} else {
		Queue = memory.NewMemoryProvider()
	}

	return Queue
}

func StartWorker(queue types.Queue[types.Mail]) {
	var ctx context.Context

	ctx, cancel = context.WithCancel(context.Background())
	go queue.Read(ctx)
}

func StopWorker() {
	cancel()
}
