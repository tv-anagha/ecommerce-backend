package events

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/tv-anagha/ecommerce-backend/order-service/internal/model"
)

// Publisher sends domain events to Kafka.
type Publisher interface {
	PublishOrderPlaced(ctx context.Context, order *model.Order) error
}

type kafkaPublisher struct {
	producer messageProducer
}

// messageProducer is the subset of kafka.Producer used by Publisher (for testing).
type messageProducer interface {
	Publish(ctx context.Context, key string, value []byte) error
}

// NewPublisher creates a Kafka-backed event publisher.
func NewPublisher(producer messageProducer) Publisher {
	return &kafkaPublisher{producer: producer}
}

// PublishOrderPlaced serializes an order.placed event and writes it to Kafka.
// The message key is userId so all events for one user land in the same partition.
func (p *kafkaPublisher) PublishOrderPlaced(ctx context.Context, order *model.Order) error {
	event := NewOrderPlacedEvent(order)
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal order.placed: %w", err)
	}

	key := strconv.FormatUint(uint64(order.UserID), 10)
	if err := p.producer.Publish(ctx, key, payload); err != nil {
		return fmt.Errorf("publish order.placed: %w", err)
	}
	return nil
}
