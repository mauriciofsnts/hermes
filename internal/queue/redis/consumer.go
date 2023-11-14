package redis

import (
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

type Consumer[T any] struct {
	Client *redis.Client
	Topic  string
}

func NewConsumer[T any](client *redis.Client, topic string) *Consumer[T] {
	return &Consumer[T]{
		Client: client,
		Topic:  topic,
	}
}

type ReadData[T any] struct {
	Data *T
	Err  error
}

func (c *Consumer[T]) Read(ch chan<- ReadData[T]) {
	pubsub := c.Client.Subscribe(ctx, c.Topic)

	for {
		msg, err := pubsub.ReceiveMessage(ctx)

		if err != nil {
			ch <- ReadData[T]{nil, err}
			continue
		}

		var model T

		err = json.Unmarshal([]byte(msg.Payload), &model)

		if err != nil {
			ch <- ReadData[T]{nil, err}
			continue
		}

		ch <- ReadData[T]{&model, nil}
	}

}

func (c *Consumer[T]) Close() error {
	return c.Client.Close()
}
