package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/wallet-service/internal/config"
)

type BlockCypherProvider struct {
	token   string
	baseURL string
	client  *http.Client
}

func NewBlockCypherProvider(cfg *config.Config) *BlockCypherProvider {
	token := cfg.CryptoAPIKeys["BLOCKCYPHER_TOKEN"]
	return &BlockCypherProvider{
		token:   token,
		baseURL: "https://api.blockcypher.com/v1",
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (p *BlockCypherProvider) Name() string {
	return "BlockCypher"
}

func (p *BlockCypherProvider) GetBalance(currency, address string) (float64, error) {
	if currency != "BTC" {
		return 0, fmt.Errorf("blockcypher only supports BTC (in this implementation)")
	}

	// https://api.blockcypher.com/v1/btc/main/addrs/ADDR/balance
	url := fmt.Sprintf("%s/btc/main/addrs/%s/balance?token=%s", p.baseURL, address, p.token)

	resp, err := p.client.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return 0, fmt.Errorf("blockcypher error %d", resp.StatusCode)
	}

	type BlockCypherBalance struct {
		Balance      int64 `json:"balance"` // Satoshis
		Unconfirmed  int64 `json:"unconfirmed_balance"`
		FinalBalance int64 `json:"final_balance"`
	}

	var bal BlockCypherBalance
	if err := json.NewDecoder(resp.Body).Decode(&bal); err != nil {
		return 0, err
	}

	// Convert Satoshis to BTC
	return float64(bal.FinalBalance) / 100000000.0, nil
}

func (p *BlockCypherProvider) BroadcastTransaction(currency, txHex string) (string, error) {
	if currency != "BTC" {
		return "", fmt.Errorf("blockcypher only supports BTC")
	}

	url := fmt.Sprintf("%s/btc/main/txs/push?token=%s", p.baseURL, p.token)

	reqBody := map[string]string{"tx": txHex}
	jsonData, _ := json.Marshal(reqBody)

	resp, err := p.client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("blockcypher broadcast error %d: %s", resp.StatusCode, string(bodyBytes))
	}

	type BlockCypherPush struct {
		Tx struct {
			Hash string `json:"hash"`
		} `json:"tx"`
	}

	var res BlockCypherPush
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}

	return res.Tx.Hash, nil
}
