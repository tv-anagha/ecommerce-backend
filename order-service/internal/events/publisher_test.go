package events

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/tv-anagha/ecommerce-backend/order-service/internal/model"
)

type mockProducer struct {
	key     string
	payload []byte
	err     error
}

func (m *mockProducer) Publish(_ context.Context, key string, value []byte) error {
	m.key = key
	m.payload = value
	return m.err
}

func TestPublishOrderPlaced(t *testing.T) {
	producer := &mockProducer{}
	pub := NewPublisher(producer)

	order := &model.Order{
		ID:          42,
		UserID:      7,
		TotalAmount: 199.99,
		Items: []model.OrderItem{
			{ProductID: 1, ProductName: "Widget", Price: 99.99, Quantity: 2},
		},
	}

	if err := pub.PublishOrderPlaced(context.Background(), order); err != nil {
		t.Fatalf("PublishOrderPlaced: %v", err)
	}

	if producer.key != "7" {
		t.Fatalf("key = %q, want %q", producer.key, "7")
	}

	var event OrderPlacedEvent
	if err := json.Unmarshal(producer.payload, &event); err != nil {
		t.Fatalf("unmarshal payload: %v", err)
	}
	if event.EventType != TypeOrderPlaced {
		t.Fatalf("eventType = %q, want %q", event.EventType, TypeOrderPlaced)
	}
	if event.OrderID != 42 || event.UserID != 7 {
		t.Fatalf("unexpected ids: orderId=%d userId=%d", event.OrderID, event.UserID)
	}
	if len(event.Items) != 1 || event.Items[0].ProductName != "Widget" {
		t.Fatalf("unexpected items: %+v", event.Items)
	}
}

func TestNoopPublisher(t *testing.T) {
	pub := NewNoopPublisher()
	order := &model.Order{ID: 1, UserID: 2}
	if err := pub.PublishOrderPlaced(context.Background(), order); err != nil {
		t.Fatalf("noop publish should not fail: %v", err)
	}
}

func TestKafkaConfigured(t *testing.T) {
	t.Setenv("KAFKA_ENABLED", "")
	t.Setenv("KAFKA_BROKERS", "")
	if KafkaConfigured() {
		t.Fatal("expected disabled when KAFKA_BROKERS is empty")
	}

	t.Setenv("KAFKA_BROKERS", "localhost:9092")
	if !KafkaConfigured() {
		t.Fatal("expected enabled when KAFKA_BROKERS is set")
	}

	t.Setenv("KAFKA_ENABLED", "false")
	if KafkaConfigured() {
		t.Fatal("expected disabled when KAFKA_ENABLED=false")
	}
}
