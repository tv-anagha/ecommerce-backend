package repository

import (
	"errors"

	"github.com/tv-anagha/ecommerce-backend/order-service/internal/database"
	"github.com/tv-anagha/ecommerce-backend/order-service/internal/model"
	"gorm.io/gorm"
)

var ErrOrderNotFound = errors.New("order not found")

type OrderRepository struct{}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{}
}

func (r *OrderRepository) Create(order *model.Order) error {
	return database.DB.Create(order).Error
}

func (r *OrderRepository) FindByID(id uint) (*model.Order, error) {
	var order model.Order
	err := database.DB.Preload("Items").First(&order, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrOrderNotFound
	}
	if err != nil {
		return nil, err
	}
	return &order, nil
}
