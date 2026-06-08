package repository

import (
	"errors"

	"github.com/tv-anagha/ecommerce-backend/order-service/internal/database"
	"github.com/tv-anagha/ecommerce-backend/order-service/internal/model"
	"gorm.io/gorm"
)

var ErrOrderNotFound = errors.New("order not found")

// OrderRepository is the repository for the order model.
type OrderRepository struct{}

// NewOrderRepository creates a new OrderRepository.
func NewOrderRepository() *OrderRepository {
	// Return a new OrderRepository.
	return &OrderRepository{}
}

func (r *OrderRepository) Create(order *model.Order) error {
	// Create a new order in the database.
	return database.DB.Create(order).Error
}

func (r *OrderRepository) FindByID(id uint) (*model.Order, error) {
	// Find an order by its ID.
	var order model.Order
	// Preload the items of the order.
	err := database.DB.Preload("Items").First(&order, id).Error
	// If the order is not found, return an error.
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrOrderNotFound
	}
	// If there is an error, return an error.
	if err != nil {
		return nil, err
	}
	// Return the order.
	return &order, nil
}
