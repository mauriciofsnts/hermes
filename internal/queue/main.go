package queue

import (
	"context"

	"github.com/mauriciofsnts/hermes/internal/config"
	"github.com/mauriciofsnts/hermes/internal/queue/kafka"
	"github.com/mauriciofsnts/hermes/internal/queue/memory"
	"github.com/mauriciofsnts/hermes/internal/queue/redis"
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

		Queue = kafka.NewKafkaQueue()
	} else if redisEnabled {
		Queue = redis.NewRedisQueue()
	} else {
		Queue = memory.NewMemoryQueue()
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
