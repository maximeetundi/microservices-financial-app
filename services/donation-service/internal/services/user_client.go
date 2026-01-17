package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type UserClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewUserClient() *UserClient {
	baseURL := os.Getenv("AUTH_SERVICE_URL")
	if baseURL == "" {
		baseURL = "http://auth-service:8081"
	}
	return &UserClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// VerifyPin calls auth-service to verify the PIN
func (c *UserClient) VerifyPin(userID, pin, token string) error {
	url := fmt.Sprintf("%s/api/v1/internal/users/pin/verify", c.baseURL)

	reqBody := map[string]string{
		"pin": pin,
	}
	jsonData, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// We need to pass the user_id header so auth-service knows who we are verifying
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", userID)
	req.Header.Set("Authorization", token) // Forward the auth token

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call auth service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	// Read error message
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusBadRequest {
		return fmt.Errorf("PIN verification failed: %s", string(body))
	}

	return fmt.Errorf("auth service returned status %d", resp.StatusCode)
}
