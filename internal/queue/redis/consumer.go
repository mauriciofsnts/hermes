package redis

import (
	"encoding/json"

	"github.com/mauriciofsnts/hermes/internal/types"
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

func (c *Consumer[T]) Read(ch chan<- types.ReadData[T]) {
	pubsub := c.Client.Subscribe(ctx, c.Topic)

	for {
		msg, err := pubsub.ReceiveMessage(ctx)

		if err != nil {
			ch <- types.ReadData[T]{Data: nil, Err: err}
			continue
		}

		var model T

		err = json.Unmarshal([]byte(msg.Payload), &model)

		if err != nil {
			ch <- types.ReadData[T]{Data: nil, Err: err}
			continue
		}

		ch <- types.ReadData[T]{Data: &model, Err: nil}
	}

}

func (c *Consumer[T]) Close() error {
	return c.Client.Close()
}
