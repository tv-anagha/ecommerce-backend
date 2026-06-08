package service

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/tv-anagha/ecommerce-backend/order-service/internal/client"
	"github.com/tv-anagha/ecommerce-backend/order-service/internal/events"
	"github.com/tv-anagha/ecommerce-backend/order-service/internal/model"
	"github.com/tv-anagha/ecommerce-backend/order-service/internal/repository"
)

var ErrEmptyCart = errors.New("cart is empty")

type OrderService struct {
	repo          *repository.OrderRepository
	cartClient    *client.CartClient
	productClient *client.ProductClient
	publisher     events.Publisher
}

func NewOrderService(	
	repo *repository.OrderRepository,
	cartClient *client.CartClient,
	productClient *client.ProductClient,
	publisher events.Publisher,
) *OrderService {
	return &OrderService{
		repo:          repo,
		cartClient:    cartClient,
		productClient: productClient,
		publisher:     publisher,
	}
}

func (s *OrderService) PlaceOrder(ctx context.Context, userID uint) (*model.Order, error) {
	cartItems, err := s.cartClient.GetCart(userID)
	if err != nil {
		return nil, err
	}
	if len(cartItems) == 0 {
		return nil, ErrEmptyCart
	}

	orderItems := make([]model.OrderItem, 0, len(cartItems))
	var total float64

	for _, item := range cartItems {
		product, err := s.productClient.GetProduct(item.ProductID)
		if err != nil {
			return nil, err
		}

		lineTotal := product.Price * float64(item.Quantity)
		total += lineTotal

		orderItems = append(orderItems, model.OrderItem{
			ProductID:   item.ProductID,
			ProductName: product.Name,
			Price:       product.Price,
			Quantity:    item.Quantity,
		})
	}

	order := &model.Order{
		UserID:      userID,
		TotalAmount: total,
		Items:       orderItems,
	}
	if err := s.repo.Create(order); err != nil {
		return nil, err
	}

	for _, item := range cartItems {
		if err := s.cartClient.RemoveItem(userID, item.ProductID); err != nil {
			log.Printf("order-service: failed to clear cart item userId=%d productId=%d: %v",
				userID, item.ProductID, err)
		}
	}

	// Publish AFTER database commit so we never emit events for failed orders.
	// If publish fails, log error but still return 201 — notification is best-effort in Phase 1.
	if err := s.publisher.PublishOrderPlaced(ctx, order); err != nil {
		log.Printf("kafka: failed to publish order.placed for order %d: %v", order.ID, err)
	}

	return order, nil
}

func (s *OrderService) GetOrder(id uint) (*model.Order, error) {
	order, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (s *OrderService) CheckoutErrorMessage(err error) string {
	switch {
	case errors.Is(err, ErrEmptyCart):
		return err.Error()
	case errors.Is(err, client.ErrProductNotFound):
		return "cart contains a product that no longer exists"
	default:
		return fmt.Sprintf("checkout failed: %v", err)
	}
}
