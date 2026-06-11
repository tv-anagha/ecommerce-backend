package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	kafkago "github.com/segmentio/kafka-go"
)

const orderPlacedEventType = "order.placed"

// OrderItemDTO mirrors a line item in the order.placed event payload.
type OrderItemDTO struct {
	ProductID   uint    `json:"productId"`
	ProductName string  `json:"productName"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
}

// OrderPlacedEvent mirrors the JSON payload published by order-service.
// Each consumer group receives the same bytes; this struct must match order-service's JSON.
type OrderPlacedEvent struct {
	EventType   string         `json:"eventType"`
	OrderID     uint           `json:"orderId"`
	UserID      uint           `json:"userId"`
	TotalAmount float64        `json:"totalAmount"`
	Items       []OrderItemDTO `json:"items"`
	PlacedAt    time.Time      `json:"placedAt"`
}

var (
	ErrInvalidPayload = errors.New("invalid order.placed payload")
	ErrUnexpectedType = errors.New("unexpected event type")
)

// OrderConsumer reads order.placed messages from Kafka and logs fulfillment actions.
type OrderConsumer struct {
	reader *kafkago.Reader
}

// NewOrderConsumer builds a consumer group reader for the order.placed topic.
// GroupID is the critical setting: Kafka tracks offsets per group, so a new group
// name means this service gets its own copy of every message on the topic —
// independent of notification-service or any other consumer.
func NewOrderConsumer() *OrderConsumer {
	brokers := strings.Split(env("KAFKA_BROKERS", "localhost:9092"), ",")
	topic := env("KAFKA_TOPIC", "order.placed")
	groupID := env("KAFKA_GROUP_ID", "fulfillment-service")

	log.Printf("[kafka] connecting broker=%s topic=%s group=%s", strings.Join(brokers, ","), topic, groupID)

	return &OrderConsumer{
		reader: kafkago.NewReader(kafkago.ReaderConfig{
			Brokers:        brokers,
			Topic:          topic,
			GroupID:        groupID,
			MinBytes:       1,
			MaxBytes:       10e6,
			CommitInterval: time.Second,
			StartOffset:    kafkago.LastOffset,
		}),
	}
}

// Run blocks until ctx is cancelled, processing messages as they arrive.
func (c *OrderConsumer) Run(ctx context.Context) error {
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return fmt.Errorf("read message: %w", err)
		}

		if err := c.handleMessage(msg); err != nil {
			if errors.Is(err, ErrInvalidPayload) || errors.Is(err, ErrUnexpectedType) {
				log.Printf("[kafka] skipping message partition=%d offset=%d: %v", msg.Partition, msg.Offset, err)
				continue
			}
			return err
		}
	}
}

// handleMessage runs only inside fulfillment-service — no imports from
// notification-service or analytics-service. Both services react to the same
// JSON event independently via their own Kafka consumer groups.
func (c *OrderConsumer) handleMessage(msg kafkago.Message) error {
	event, err := parseOrderPlaced(msg.Value)
	if err != nil {
		return err
	}

	log.Printf("[fulfillment] order.placed received — orderId=%d userId=%d items=%d",
		event.OrderID, event.UserID, len(event.Items))

	if len(event.Items) > 0 {
		picks := make([]string, 0, len(event.Items))
		for _, item := range event.Items {
			picks = append(picks, fmt.Sprintf("%dx %s", item.Quantity, item.ProductName))
		}
		log.Printf("[fulfillment] picking: %s", strings.Join(picks, ", "))
	}

	log.Printf("[fulfillment] shipment queued for order #%d", event.OrderID)
	return nil
}

func parseOrderPlaced(value []byte) (OrderPlacedEvent, error) {
	var event OrderPlacedEvent
	if err := json.Unmarshal(value, &event); err != nil {
		return OrderPlacedEvent{}, fmt.Errorf("%w: %v", ErrInvalidPayload, err)
	}
	if event.EventType != "" && event.EventType != orderPlacedEventType {
		return OrderPlacedEvent{}, fmt.Errorf("%w: %q", ErrUnexpectedType, event.EventType)
	}
	if event.OrderID == 0 {
		return OrderPlacedEvent{}, fmt.Errorf("%w: missing orderId", ErrInvalidPayload)
	}
	return event, nil
}

// Close releases the Kafka reader.
func (c *OrderConsumer) Close() error {
	return c.reader.Close()
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
