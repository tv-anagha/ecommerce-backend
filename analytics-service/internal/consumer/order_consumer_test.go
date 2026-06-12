package consumer

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestParseOrderPlacedValid(t *testing.T) {
	raw, err := json.Marshal(OrderPlacedEvent{
		EventType:   "order.placed",
		OrderID:     10,
		UserID:      3,
		TotalAmount: 50,
		Items:       []OrderItemDTO{{ProductID: 1, ProductName: "Widget", Quantity: 2}},
	})
	if err != nil {
		t.Fatal(err)
	}

	event, err := parseOrderPlaced(raw)
	if err != nil {
		t.Fatalf("parseOrderPlaced: %v", err)
	}
	if event.OrderID != 10 || len(event.Items) != 1 {
		t.Fatalf("unexpected event: %+v", event)
	}
}

func TestParseOrderPlacedRejectsBadJSON(t *testing.T) {
	_, err := parseOrderPlaced([]byte("{not-json"))
	if !errors.Is(err, ErrInvalidPayload) {
		t.Fatalf("expected ErrInvalidPayload, got %v", err)
	}
}

func TestParseOrderPlacedRejectsWrongType(t *testing.T) {
	raw, err := json.Marshal(OrderPlacedEvent{
		EventType: "order.cancelled",
		OrderID:   1,
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = parseOrderPlaced(raw)
	if !errors.Is(err, ErrUnexpectedType) {
		t.Fatalf("expected ErrUnexpectedType, got %v", err)
	}
}
