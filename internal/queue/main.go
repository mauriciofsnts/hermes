package queue

import (
	"context"

	"github.com/mauriciofsnts/hermes/internal/config"
	"github.com/mauriciofsnts/hermes/internal/queue/kafka"
	"github.com/mauriciofsnts/hermes/internal/queue/memory"
	"github.com/mauriciofsnts/hermes/internal/queue/redis"
	"github.com/mauriciofsnts/hermes/internal/types"
)

var queue types.Queue[types.Email]
var cancel context.CancelFunc

func NewQueue() types.Queue[types.Email] {
	kafkaEnabled := config.Hermes.Kafka.Enabled
	redisEnabled := config.Hermes.Redis.Enabled

	if kafkaEnabled {
		err := kafka.CreateTopic()

		if err != nil {
			panic(err)
		}

		queue = kafka.NewKafkaQueue()
	} else if redisEnabled {
		queue = redis.NewRedisQueue()
	} else {
		queue = memory.NewMemoryQueue()
	}

	return queue
}

func StartWorker(queue types.Queue[types.Email]) {
	var ctx context.Context

	ctx, cancel = context.WithCancel(context.Background())
	go queue.Read(ctx)
}

func StopWorker() {
	cancel()
}
