package redis

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

const (
	Topic = "hermes"
)

type Consumer[T any] struct {
	Client *redis.Client
	Ctx    context.Context
	Topic  string
}

func NewConsumer[T any](client *redis.Client, topic string) *Consumer[T] {
	return &Consumer[T]{
		Client: client,
		Topic:  topic,
	}
}

func (c *Consumer[T]) Read(callback func(*T, error)) error {
	pubsub := c.Client.Subscribe(c.Ctx, Topic)

	for {
		msg, err := pubsub.ReceiveMessage(c.Ctx)

		if err != nil {
			callback(nil, err)
			continue
		}

		var model T

		err = json.Unmarshal([]byte(msg.Payload), &model)

		if err != nil {
			callback(nil, err)
			continue
		}

		callback(&model, nil)
	}

}
