package service

import (
	"github.com/tv-anagha/ecommerce-backend/product-service/internal/model"
	"github.com/tv-anagha/ecommerce-backend/product-service/internal/repository"
)

type ProductService struct {
	repo *repository.ProductRepository
}

func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) ListProducts() ([]model.Product, error) {
	return s.repo.List()
}

func (s *ProductService) GetProduct(id int) (model.Product, error) {
	return s.repo.GetByID(id)
}
