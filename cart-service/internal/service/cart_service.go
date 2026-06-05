package service

import (
	"errors"

	"github.com/tv-anagha/ecommerce-backend/cart-service/internal/client"
	"github.com/tv-anagha/ecommerce-backend/cart-service/internal/model"
	"github.com/tv-anagha/ecommerce-backend/cart-service/internal/repository"
)

var ErrInvalidQuantity = errors.New("quantity must be greater than zero")

type CartService struct {
	repo          *repository.CartRepository
	productClient *client.ProductClient
}

func NewCartService(repo *repository.CartRepository, productClient *client.ProductClient) *CartService {
	return &CartService{
		repo:          repo,
		productClient: productClient,
	}
}

func (s *CartService) GetCart(userID uint) ([]model.CartItem, error) {
	return s.repo.ListByUserID(userID)
}

func (s *CartService) AddItem(userID, productID uint, quantity int) (model.CartItem, error) {
	if quantity <= 0 {
		return model.CartItem{}, ErrInvalidQuantity
	}
	if err := s.productClient.ProductExists(productID); err != nil {
		return model.CartItem{}, err
	}

	existing, err := s.repo.GetByUserAndProduct(userID, productID)
	if errors.Is(err, repository.ErrCartItemNotFound) {
		return s.repo.UpsertItem(userID, productID, quantity)
	}
	if err != nil {
		return model.CartItem{}, err
	}

	return s.repo.UpsertItem(userID, productID, existing.Quantity+quantity)
}

func (s *CartService) UpdateItem(userID, productID uint, quantity int) (model.CartItem, error) {
	if quantity <= 0 {
		return model.CartItem{}, ErrInvalidQuantity
	}
	return s.repo.UpdateQuantity(userID, productID, quantity)
}

func (s *CartService) RemoveItem(userID, productID uint) error {
	return s.repo.RemoveItem(userID, productID)
}
