package queue

import (
	"github.com/mauriciofsnts/hermes/internal/config"
	"github.com/mauriciofsnts/hermes/internal/queue/kafka"
	"github.com/mauriciofsnts/hermes/internal/queue/memory"
	"github.com/mauriciofsnts/hermes/internal/queue/redis"
	"github.com/mauriciofsnts/hermes/internal/types"
)

var queue types.Queue[types.Email]

func NewQueue() types.Queue[types.Email] {
	kafkaEnabled := config.Hermes.Kafka.Enabled
	redisEnabled := config.Hermes.Redis.Enabled

	if kafkaEnabled {
		err := kafka.CreateTopic()

		if err != nil {
			panic(err)
		}

		queue = kafka.NewKafkaQueue()
		return queue
	} else if redisEnabled {
		queue = redis.NewRedisQueue()
		return queue
	} else {
		queue = memory.NewMemoryQueue()
		return queue
	}
}
