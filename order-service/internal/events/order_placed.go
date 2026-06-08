package events

import (
	"time"

	"github.com/tv-anagha/ecommerce-backend/order-service/internal/model"
)

const TypeOrderPlaced = "order.placed"

// OrderItemDTO is a line item included in the order.placed event payload.
type OrderItemDTO struct {
	ProductID   uint    `json:"productId"`
	ProductName string  `json:"productName"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
}

// OrderPlacedEvent is published to Kafka after an order is saved.
type OrderPlacedEvent struct {
	EventType   string         `json:"eventType"`
	OrderID     uint           `json:"orderId"`
	UserID      uint           `json:"userId"`
	TotalAmount float64        `json:"totalAmount"`
	Items       []OrderItemDTO `json:"items"`
	PlacedAt    time.Time      `json:"placedAt"`
}

// NewOrderPlacedEvent builds the event payload from a persisted order.
func NewOrderPlacedEvent(order *model.Order) OrderPlacedEvent {
	items := make([]OrderItemDTO, 0, len(order.Items))
	for _, item := range order.Items {
		items = append(items, OrderItemDTO{
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Price:       item.Price,
			Quantity:    item.Quantity,
		})
	}

	placedAt := order.CreatedAt
	if placedAt.IsZero() {
		placedAt = time.Now().UTC()
	}

	return OrderPlacedEvent{
		EventType:   TypeOrderPlaced,
		OrderID:     order.ID,
		UserID:      order.UserID,
		TotalAmount: order.TotalAmount,
		Items:       items,
		PlacedAt:    placedAt,
	}
}
