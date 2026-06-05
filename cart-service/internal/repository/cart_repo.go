package repository

import (
	"errors"

	"github.com/tv-anagha/ecommerce-backend/cart-service/internal/database"
	"github.com/tv-anagha/ecommerce-backend/cart-service/internal/model"
	"gorm.io/gorm"
)

var ErrCartItemNotFound = errors.New("cart item not found")

type CartRepository struct{}

func NewCartRepository() *CartRepository {
	return &CartRepository{}
}

func (r *CartRepository) ListByUserID(userID uint) ([]model.CartItem, error) {
	var items []model.CartItem
	if err := database.DB.Where("user_id = ?", userID).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *CartRepository) GetByUserAndProduct(userID, productID uint) (model.CartItem, error) {
	var item model.CartItem
	err := database.DB.Where("user_id = ? AND product_id = ?", userID, productID).First(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return model.CartItem{}, ErrCartItemNotFound
	}
	return item, err
}

func (r *CartRepository) UpsertItem(userID, productID uint, quantity int) (model.CartItem, error) {
	item, err := r.GetByUserAndProduct(userID, productID)
	if errors.Is(err, ErrCartItemNotFound) {
		item = model.NewCartItem(userID, productID, quantity)
		if err := database.DB.Create(&item).Error; err != nil {
			return model.CartItem{}, err
		}
		return item, nil
	}
	if err != nil {
		return model.CartItem{}, err
	}

	item.Quantity = quantity
	if err := database.DB.Save(&item).Error; err != nil {
		return model.CartItem{}, err
	}
	return item, nil
}

func (r *CartRepository) UpdateQuantity(userID, productID uint, quantity int) (model.CartItem, error) {
	item, err := r.GetByUserAndProduct(userID, productID)
	if err != nil {
		return model.CartItem{}, err
	}

	item.Quantity = quantity
	if err := database.DB.Save(&item).Error; err != nil {
		return model.CartItem{}, err
	}
	return item, nil
}

func (r *CartRepository) RemoveItem(userID, productID uint) error {
	result := database.DB.Where("user_id = ? AND product_id = ?", userID, productID).Delete(&model.CartItem{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrCartItemNotFound
	}
	return nil
}
