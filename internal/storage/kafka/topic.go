package kafka

import (
	"net"

	"github.com/pauloo27/logger"
	kafkaGo "github.com/segmentio/kafka-go"
)

const (
	// TODO: Move to config
	EmailTopic = "kafka-email-topic"
)

var defaultTopics = []string{
	EmailTopic,
}

func CreateTopic() error {
	connection, err := kafkaGo.Dial("tcp", net.JoinHostPort("localhost", "9092"))

	if err != nil {
		logger.Error("Failed to connect to Kafka", err)
		return err
	}

	topics := make([]kafkaGo.TopicConfig, len(defaultTopics))

	for i, topic := range defaultTopics {
		topics[i] = NewTopic(topic)
	}

	err = connection.CreateTopics(topics...)

	if err != nil {
		logger.Error("Failed to create topic", err)
		return err
	}

	return nil
}

func NewTopic(topicName string) kafkaGo.TopicConfig {
	return kafkaGo.TopicConfig{
		Topic:             topicName,
		NumPartitions:     1,
		ReplicationFactor: 1,
	}
}
