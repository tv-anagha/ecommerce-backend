package kafka

import (
	"context"
	"os"
	"strings"
	"time"

	kafkago "github.com/segmentio/kafka-go"
)

// Producer wraps a kafka-go Writer for publishing JSON events.
type Producer struct {
	writer *kafkago.Writer
	topic  string
}

// NewProducer connects to the Kafka broker(s) from KAFKA_BROKERS.
func NewProducer() *Producer {
	brokers := strings.Split(env("KAFKA_BROKERS", "localhost:9092"), ",")
	topic := env("KAFKA_TOPIC", "order.placed")

	return &Producer{
		topic: topic,
		writer: &kafkago.Writer{
			Addr:         kafkago.TCP(brokers...),
			Topic:        topic,
			Balancer:     &kafkago.LeastBytes{},
			RequiredAcks: kafkago.RequireOne,
			// Short timeout for Phase 1; checkout still succeeds if publish fails.
			WriteTimeout: 10 * time.Second,
		},
	}
}

// Publish sends a message to the configured topic.
// key is used for partition routing (e.g. userId keeps one user's events ordered).
func (p *Producer) Publish(ctx context.Context, key string, value []byte) error {
	return p.writer.WriteMessages(ctx, kafkago.Message{
		Key:   []byte(key),
		Value: value,
	})
}

// Close releases the underlying writer.
func (p *Producer) Close() error {
	return p.writer.Close()
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
