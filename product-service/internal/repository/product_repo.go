package repository

import (
	"errors"

	"github.com/tv-anagha/ecommerce-backend/product-service/internal/database"
	"github.com/tv-anagha/ecommerce-backend/product-service/internal/model"
	"gorm.io/gorm"
)

var ErrProductNotFound = errors.New("product not found")

type ProductRepository struct{}

func NewProductRepository() *ProductRepository {
	return &ProductRepository{}
}

func (r *ProductRepository) List() ([]model.Product, error) {
	var products []model.Product
	if err := database.DB.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductRepository) GetByID(id int) (model.Product, error) {
	var product model.Product
	err := database.DB.First(&product, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return model.Product{}, ErrProductNotFound
	}
	return product, err
}
