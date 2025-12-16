package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/api-gateway/internal/config"
)

type ServiceManager struct {
	config *config.Config
	client *http.Client
}

func NewServiceManager(cfg *config.Config) *ServiceManager {
	return &ServiceManager{
		config: cfg,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

type ServiceResponse struct {
	StatusCode int
	Body       []byte
	Headers    http.Header
}

func (sm *ServiceManager) CallService(ctx context.Context, serviceName, method, endpoint string, body interface{}, headers map[string]string) (*ServiceResponse, error) {
	serviceURL := sm.getServiceURL(serviceName)
	if serviceURL == "" {
		return nil, fmt.Errorf("unknown service: %s", serviceName)
	}

	url := serviceURL + endpoint
	
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := sm.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call service: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return &ServiceResponse{
		StatusCode: resp.StatusCode,
		Body:       responseBody,
		Headers:    resp.Header,
	}, nil
}

func (sm *ServiceManager) getServiceURL(serviceName string) string {
	switch serviceName {
	case "auth":
		return sm.config.AuthServiceURL
	case "user":
		return sm.config.UserServiceURL
	case "wallet":
		return sm.config.WalletServiceURL
	case "transfer":
		return sm.config.TransferServiceURL
	case "exchange":
		return sm.config.ExchangeServiceURL
	case "card":
		return sm.config.CardServiceURL
	case "notification":
		return sm.config.NotificationServiceURL
	default:
		return ""
	}
}

// Auth Service calls
func (sm *ServiceManager) Login(ctx context.Context, email, password string) (*ServiceResponse, error) {
	body := map[string]string{
		"email":    email,
		"password": password,
	}
	return sm.CallService(ctx, "auth", "POST", "/api/v1/login", body, nil)
}

func (sm *ServiceManager) Register(ctx context.Context, userData map[string]interface{}) (*ServiceResponse, error) {
	return sm.CallService(ctx, "auth", "POST", "/api/v1/register", userData, nil)
}

func (sm *ServiceManager) RefreshToken(ctx context.Context, refreshToken string) (*ServiceResponse, error) {
	body := map[string]string{
		"refresh_token": refreshToken,
	}
	return sm.CallService(ctx, "auth", "POST", "/api/v1/refresh", body, nil)
}

// User Service calls
func (sm *ServiceManager) GetUser(ctx context.Context, userID string, authToken string) (*ServiceResponse, error) {
	headers := map[string]string{
		"Authorization": "Bearer " + authToken,
	}
	return sm.CallService(ctx, "user", "GET", "/users/"+userID, nil, headers)
}

func (sm *ServiceManager) UpdateUser(ctx context.Context, userID string, userData map[string]interface{}, authToken string) (*ServiceResponse, error) {
	headers := map[string]string{
		"Authorization": "Bearer " + authToken,
	}
	return sm.CallService(ctx, "user", "PUT", "/users/"+userID, userData, headers)
}

// Wallet Service calls
func (sm *ServiceManager) GetWallets(ctx context.Context, userID string, authToken string) (*ServiceResponse, error) {
	headers := map[string]string{
		"Authorization": "Bearer " + authToken,
	}
	return sm.CallService(ctx, "wallet", "GET", "/wallets?user_id="+userID, nil, headers)
}

func (sm *ServiceManager) CreateWallet(ctx context.Context, walletData map[string]interface{}, authToken string) (*ServiceResponse, error) {
	headers := map[string]string{
		"Authorization": "Bearer " + authToken,
	}
	return sm.CallService(ctx, "wallet", "POST", "/wallets", walletData, headers)
}

// Transfer Service calls
func (sm *ServiceManager) CreateTransfer(ctx context.Context, transferData map[string]interface{}, authToken string) (*ServiceResponse, error) {
	headers := map[string]string{
		"Authorization": "Bearer " + authToken,
	}
	return sm.CallService(ctx, "transfer", "POST", "/transfers", transferData, headers)
}

func (sm *ServiceManager) GetTransferHistory(ctx context.Context, userID string, authToken string) (*ServiceResponse, error) {
	headers := map[string]string{
		"Authorization": "Bearer " + authToken,
	}
	return sm.CallService(ctx, "transfer", "GET", "/transfers?user_id="+userID, nil, headers)
}

// Exchange Service calls
func (sm *ServiceManager) GetExchangeRates(ctx context.Context) (*ServiceResponse, error) {
	return sm.CallService(ctx, "exchange", "GET", "/rates", nil, nil)
}

func (sm *ServiceManager) CreateExchange(ctx context.Context, exchangeData map[string]interface{}, authToken string) (*ServiceResponse, error) {
	headers := map[string]string{
		"Authorization": "Bearer " + authToken,
	}
	return sm.CallService(ctx, "exchange", "POST", "/exchange", exchangeData, headers)
}