package secrets

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/hashicorp/vault/api"
)

type VaultClient struct {
	client *api.Client
}

func NewVaultClient() (*VaultClient, error) {
	config := api.DefaultConfig()

	// Read from env if set, otherwise default
	if addr := os.Getenv("VAULT_ADDR"); addr != "" {
		config.Address = addr
	}

	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	// In a real prod environment, we would use AppRole or Kubernetes Auth
	// For this setup, we rely on VAULT_TOKEN being present (injected by init script or env)
	if token := os.Getenv("VAULT_TOKEN"); token != "" {
		client.SetToken(token)
	}

	return &VaultClient{client: client}, nil
}

func (v *VaultClient) GetSecret(path string) (map[string]interface{}, error) {
	// Add retry logic
	var err error
	for i := 0; i < 3; i++ {
		secret, err := v.client.Logical().Read(path)
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		if secret == nil {
			return nil, errors.New("secret not found")
		}
		return secret.Data, nil
	}
	return nil, fmt.Errorf("failed to read secret after retries: %v", err)
}

func (v *VaultClient) WriteSecret(path string, data map[string]interface{}) error {
	_, err := v.client.Logical().Write(path, data)
	return err
}
