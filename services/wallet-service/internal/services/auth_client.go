package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// AuthClient provides HTTP client to communicate with auth-service
type AuthClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewAuthClient creates a new auth client
func NewAuthClient() *AuthClient {
	authURL := os.Getenv("AUTH_SERVICE_URL")
	if authURL == "" {
		authURL = "http://auth-service:8080"
	}
	return &AuthClient{
		baseURL:    authURL,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// VerifyPinResponse represents the response from auth-service
type VerifyPinResponse struct {
	Valid       bool   `json:"valid"`
	Message     string `json:"message"`
	LockedUntil string `json:"locked_until,omitempty"`
}

// VerifyPin verifies the user's PIN via auth-service
func (c *AuthClient) VerifyPin(userID, pin, token string) (bool, error) {
	payload := map[string]string{
		"pin": pin,
	}

	jsonData, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/users/pin/verify", c.baseURL), bytes.NewBuffer(jsonData))
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", userID)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("auth service unavailable: %w", err)
	}
	defer resp.Body.Close()

	var pinResp VerifyPinResponse
	json.NewDecoder(resp.Body).Decode(&pinResp)

	if resp.StatusCode != http.StatusOK {
		if pinResp.Message != "" {
			return false, fmt.Errorf(pinResp.Message)
		}
		return false, fmt.Errorf("PIN verification failed")
	}

	return pinResp.Valid, nil
}
