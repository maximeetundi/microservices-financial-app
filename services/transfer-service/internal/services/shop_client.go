package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/config"
)

type ShopClient struct {
	baseURL    string
	httpClient *http.Client
}

type ShopDetails struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Slug        string `json:"slug"`
	LogoURL     string `json:"logo_url"`
	IsPublic    bool   `json:"is_public"`
}

func NewShopClient(cfg *config.Config) *ShopClient {
	// Assuming shop service URL is in config or default
	// We might need to add ShopServiceURL to config if not present
	// Default from config
	url := cfg.ShopServiceURL
	if url == "" {
		url = "http://shop-service:8098"
	}
	
	return &ShopClient{
		baseURL: url,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *ShopClient) GetShopByWalletID(walletID string) (*ShopDetails, error) {
	url := fmt.Sprintf("%s/api/v1/shops/by-wallet/%s", c.baseURL, walletID)
	
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get shop: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("shop not found")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("shop service returned status: %d", resp.StatusCode)
	}

	var shop ShopDetails
	if err := json.NewDecoder(resp.Body).Decode(&shop); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &shop, nil
}
