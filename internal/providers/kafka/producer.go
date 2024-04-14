package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/mauriciofsnts/hermes/internal/config"
	kafkaGo "github.com/segmentio/kafka-go"
)

type Producer[T any] struct {
	Writer *kafkaGo.Writer
	Dialer *kafkaGo.Dialer
}

func NewProducer[T any]() *Producer[T] {
	dialer := &kafkaGo.Dialer{
		Timeout:   10 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	writer := kafkaGo.NewWriter(kafkaGo.WriterConfig{
		Brokers:   config.Envs.Kafka.Brokers,
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

	return p.Writer.WriteMessages(context.Background(), kafkaGo.Message{
		Topic:  topic,
		Offset: 0,
		Key:    []byte(key),
		Value:  data,
	})
}
