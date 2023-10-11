package events

import (
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
)

type Producer[T any] struct {
	Writer *kafka.Writer
	Dialer *kafka.Dialer
}

func NewProducer[T any]() *Producer[T] {
	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:   []string{"localhost:9092"},
		Dialer:    dialer,
		BatchSize: 1,
	})

	return &Producer[T]{
		Writer: writer,
	}
}

func (p *Producer[T]) Produce(key string, value T, topic string) error {
	data, err := json.Marshal(value)

	if err != nil {
		return err
	}

	return p.Writer.WriteMessages(context.Background(), kafka.Message{
		Topic:  topic,
		Offset: 0,
		Key:    []byte(key),
		Value:  data,
	})
}
