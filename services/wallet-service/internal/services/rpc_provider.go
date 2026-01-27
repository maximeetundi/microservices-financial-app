package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/wallet-service/internal/config"
)

type RpcProvider struct {
	endpoints map[string]string // Map currency -> RPC URL
	client    *http.Client
}

func NewRpcProvider(cfg *config.Config) *RpcProvider {
	// Build endpoints map from Config or Keys
	endpoints := make(map[string]string)

	// Infura / Alchemy keys
	infuraKey := cfg.CryptoAPIKeys["INFURA"]
	alchemyKey := cfg.CryptoAPIKeys["ALCHEMY"]

	// ETH Mainnet
	if infuraKey != "" {
		endpoints["ETH"] = fmt.Sprintf("https://mainnet.infura.io/v3/%s", infuraKey)
	} else if alchemyKey != "" {
		endpoints["ETH"] = fmt.Sprintf("https://eth-mainnet.g.alchemy.com/v2/%s", alchemyKey)
	} else {
		// Fallback to public/configured RPC
		endpoints["ETH"] = cfg.BlockchainRPC["ETH"]
	}

	// BSC
	endpoints["BSC"] = cfg.BlockchainRPC["BSC"]

	// SOL (Solana uses JSON-RPC 2.0 too, but slightly different method names usually, keep simplified)
	// endpoints["SOL"] = ...

	return &RpcProvider{
		endpoints: endpoints,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (p *RpcProvider) Name() string {
	return "JSON-RPC-Provider"
}

// JSON-RPC Request/Response structures
type JsonRpcRequest struct {
	JsonRpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	Id      int           `json:"id"`
}

type JsonRpcResponse struct {
	JsonRpc string      `json:"jsonrpc"`
	Id      int         `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func (p *RpcProvider) GetBalance(currency, address string) (float64, error) {
	url, ok := p.endpoints[strings.ToUpper(currency)]
	if !ok || url == "" {
		return 0, fmt.Errorf("no RPC endpoint configured for %s", currency)
	}

	// EVM specific: eth_getBalance
	reqBody := JsonRpcRequest{
		JsonRpc: "2.0",
		Method:  "eth_getBalance",
		Params:  []interface{}{address, "latest"},
		Id:      1,
	}

	jsonData, _ := json.Marshal(reqBody)
	resp, err := p.client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return 0, fmt.Errorf("rpc error %d", resp.StatusCode)
	}

	var res JsonRpcResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return 0, err
	}

	if res.Error != nil {
		return 0, fmt.Errorf("rpc api error: %v", res.Error)
	}

	// Result is hex string (wei)
	hexBal, ok := res.Result.(string)
	if !ok {
		return 0, fmt.Errorf("invalid result format")
	}

	// Parse Hex to Int/Float
	// Remove 0x
	if strings.HasPrefix(hexBal, "0x") {
		hexBal = hexBal[2:]
	}

	val, err := strconv.ParseUint(hexBal, 16, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse balance hex: %w", err)
	}

	// Convert Wei to ETH (1e18)
	return float64(val) / 1e18, nil
}

func (p *RpcProvider) BroadcastTransaction(currency, txHex string) (string, error) {
	url, ok := p.endpoints[strings.ToUpper(currency)]
	if !ok || url == "" {
		return "", fmt.Errorf("no RPC endpoint configured for %s", currency)
	}

	// Add 0x prefix if missing for EVM
	if !strings.HasPrefix(txHex, "0x") {
		txHex = "0x" + txHex
	}

	reqBody := JsonRpcRequest{
		JsonRpc: "2.0",
		Method:  "eth_sendRawTransaction",
		Params:  []interface{}{txHex},
		Id:      1,
	}

	jsonData, _ := json.Marshal(reqBody)
	resp, err := p.client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("rpc broadcast error %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var res JsonRpcResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}

	if res.Error != nil {
		return "", fmt.Errorf("rpc api error: %v", res.Error)
	}

	txHash, ok := res.Result.(string)
	if !ok {
		return "", fmt.Errorf("invalid result format")
	}

	return txHash, nil
}
