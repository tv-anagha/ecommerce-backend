package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
)

var ErrProductNotFound = errors.New("product not found")

// Product holds pricing details fetched at checkout.
type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// ProductClient fetches product details from product-service HTTP API.
type ProductClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewProductClient() *ProductClient {
	baseURL := os.Getenv("PRODUCT_SERVICE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8081"
	}
	return &ProductClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *ProductClient) GetProduct(productID uint) (*Product, error) {
	url := fmt.Sprintf("%s/products/%d", c.baseURL, productID)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("product-service unreachable: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	switch resp.StatusCode {
	case http.StatusOK:
		var product Product
		if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
			return nil, fmt.Errorf("product-service response decode failed: %w", err)
		}
		return &product, nil
	case http.StatusNotFound:
		return nil, ErrProductNotFound
	default:
		return nil, fmt.Errorf("product-service returned status %d", resp.StatusCode)
	}
}
