package redis

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

type Consumer[T any] struct {
	Client *redis.Client
	Topic  string
	pubsub *redis.PubSub
}

type ReadData[T any] struct {
	Data *T
	Err  error
}

func NewConsumer[T any](client *redis.Client, topic string) *Consumer[T] {
	return &Consumer[T]{
		Client: client,
		Topic:  topic,
	}
}

func (c *Consumer[T]) Read(ctx context.Context, ch chan<- ReadData[T]) {
	c.pubsub = c.Client.Subscribe(ctx, c.Topic)
	defer c.pubsub.Close()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		msg, err := c.pubsub.ReceiveMessage(ctx)
		if err != nil {
			select {
			case <-ctx.Done():
				return
			case ch <- ReadData[T]{Data: nil, Err: err}:
			}
			continue
		}

		var model T
		if err := json.Unmarshal([]byte(msg.Payload), &model); err != nil {
			select {
			case <-ctx.Done():
				return
			case ch <- ReadData[T]{Data: nil, Err: err}:
			}
			continue
		}

		select {
		case <-ctx.Done():
			return
		case ch <- ReadData[T]{Data: &model, Err: nil}:
		}
	}
}

func (c *Consumer[T]) Close() error {
	if c.pubsub != nil {
		return c.pubsub.Close()
	}
	return nil
}
