package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/transfer-service/internal/config"
)

type EnterpriseClient struct {
	baseURL    string
	httpClient *http.Client
}

type EnterpriseDetails struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Status      string `json:"status"`
	Logo        string `json:"logo"`
	OwnerID     string `json:"owner_id"`
}

func NewEnterpriseClient(cfg *config.Config) *EnterpriseClient {
	// Assuming enterprise service URL is in config or default
	url := cfg.EnterpriseServiceURL
	if url == "" {
		url = "http://enterprise-service:8097"
	}
	
	return &EnterpriseClient{
		baseURL: url,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *EnterpriseClient) GetEnterprise(id string) (*EnterpriseDetails, error) {
	url := fmt.Sprintf("%s/api/v1/enterprises/%s", c.baseURL, id)
	
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get enterprise: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("enterprise not found")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("enterprise service returned status: %d", resp.StatusCode)
	}

	var ent EnterpriseDetails
	if err := json.NewDecoder(resp.Body).Decode(&ent); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &ent, nil
}
