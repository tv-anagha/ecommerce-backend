package service

import (
	"errors"

	"github.com/tv-anagha/ecommerce-backend/cart-service/internal/client"
	"github.com/tv-anagha/ecommerce-backend/cart-service/internal/model"
	"github.com/tv-anagha/ecommerce-backend/cart-service/internal/repository"
)

var ErrInvalidQuantity = errors.New("quantity must be greater than zero")
var ErrInsufficientStock = errors.New("insufficient stock available")

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

	product, err := s.productClient.GetProduct(productID)
	if err != nil {
		return model.CartItem{}, err
	}

	existing, err := s.repo.GetByUserAndProduct(userID, productID)
	var targetQuantity int
	if errors.Is(err, repository.ErrCartItemNotFound) {
		targetQuantity = quantity
	} else if err != nil {
		return model.CartItem{}, err
	} else {
		targetQuantity = existing.Quantity + quantity
	}

	if targetQuantity > product.Quantity {
		return model.CartItem{}, ErrInsufficientStock
	}

	return s.repo.UpsertItem(userID, productID, targetQuantity)
}

func (s *CartService) UpdateItem(userID, productID uint, quantity int) (model.CartItem, error) {
	if quantity <= 0 {
		return model.CartItem{}, ErrInvalidQuantity
	}

	product, err := s.productClient.GetProduct(productID)
	if err != nil {
		return model.CartItem{}, err
	}

	if quantity > product.Quantity {
		return model.CartItem{}, ErrInsufficientStock
	}

	return s.repo.UpdateQuantity(userID, productID, quantity)
}

func (s *CartService) RemoveItem(userID, productID uint) error {
	return s.repo.RemoveItem(userID, productID)
}
