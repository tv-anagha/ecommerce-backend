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

// ProductClient validates product IDs via product-service HTTP API.
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

// ProductExists returns nil when product-service responds 200 for the given id.
func (c *ProductClient) ProductExists(productID uint) error {
	url := fmt.Sprintf("%s/products/%d", c.baseURL, productID)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("product-service unreachable: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	switch resp.StatusCode {
	case http.StatusOK:
		_ = json.NewDecoder(resp.Body).Decode(&struct {
			ID uint `json:"id"`
		}{})
		return nil
	case http.StatusNotFound:
		return ErrProductNotFound
	default:
		return fmt.Errorf("product-service returned status %d", resp.StatusCode)
	}
}
