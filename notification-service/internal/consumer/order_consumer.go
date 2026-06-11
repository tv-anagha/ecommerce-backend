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

// OrderPlacedEvent mirrors the JSON payload published by order-service.
type OrderPlacedEvent struct {
	EventType   string  `json:"eventType"`
	OrderID     uint    `json:"orderId"`
	UserID      uint    `json:"userId"`
	TotalAmount float64 `json:"totalAmount"`
}

var (
	ErrInvalidPayload = errors.New("invalid order.placed payload")
	ErrUnexpectedType = errors.New("unexpected event type")
)

// OrderConsumer reads order.placed messages from Kafka and logs notifications.
type OrderConsumer struct {
	reader *kafkago.Reader
}

// NewOrderConsumer builds a consumer group reader for the order.placed topic.
func NewOrderConsumer() *OrderConsumer {
	brokers := strings.Split(env("KAFKA_BROKERS", "localhost:9092"), ",")
	topic := env("KAFKA_TOPIC", "order.placed")
	groupID := env("KAFKA_GROUP_ID", "notification-service")

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

func (c *OrderConsumer) handleMessage(msg kafkago.Message) error {
	event, err := parseOrderPlaced(msg.Value)
	if err != nil {
		return err
	}

	log.Printf("[kafka] order.placed received — orderId=%d userId=%d total=%.2f partition=%d offset=%d",
		event.OrderID, event.UserID, event.TotalAmount, msg.Partition, msg.Offset)
	log.Printf("Thank you for your order #%d! We received your purchase of $%.2f.", event.OrderID, event.TotalAmount)
	log.Println("Order received successfully.")
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
