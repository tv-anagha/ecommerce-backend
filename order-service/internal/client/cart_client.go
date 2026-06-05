package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// CartItem mirrors cart-service JSON for checkout reads.
type CartItem struct {
	ID        uint `json:"id"`
	UserID    uint `json:"userId"`
	ProductID uint `json:"productId"`
	Quantity  int  `json:"quantity"`
}

// CartClient reads and clears carts via cart-service HTTP API.
type CartClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewCartClient() *CartClient {
	baseURL := os.Getenv("CART_SERVICE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8083"
	}
	return &CartClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *CartClient) GetCart(userID uint) ([]CartItem, error) {
	url := fmt.Sprintf("%s/carts/%d", c.baseURL, userID)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("cart-service unreachable: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cart-service returned status %d", resp.StatusCode)
	}

	var items []CartItem
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return nil, fmt.Errorf("cart-service response decode failed: %w", err)
	}
	if items == nil {
		items = []CartItem{}
	}
	return items, nil
}

func (c *CartClient) RemoveItem(userID, productID uint) error {
	url := fmt.Sprintf("%s/carts/%d/items/%d", c.baseURL, userID, productID)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("cart-service unreachable: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("cart-service returned status %d", resp.StatusCode)
	}
	return nil
}
