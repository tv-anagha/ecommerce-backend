package events

import (
	"context"
	"log"

	"github.com/tv-anagha/ecommerce-backend/order-service/internal/model"
)

// noopPublisher drops events when Kafka is not configured (e.g. minimal Docker stack).
type noopPublisher struct{}

// NewNoopPublisher returns a Publisher that logs and discards events.
func NewNoopPublisher() Publisher {
	return &noopPublisher{}
}

func (p *noopPublisher) PublishOrderPlaced(_ context.Context, order *model.Order) error {
	log.Printf("kafka: disabled — skipping order.placed for order %d (set KAFKA_BROKERS to enable)", order.ID)
	return nil
}
