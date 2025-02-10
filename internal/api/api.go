package apis

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type APIResponse struct {
	Timestamp time.Time
	Source    string
	Symbol    string
	Price     float64
}

// API-specific response structures
type CoinGeckoResponse struct {
	Bitcoin struct {
		USD float64 `json:"usd"`
	} `json:"bitcoin"`
	Ethereum struct {
		USD float64 `json:"usd"`
	} `json:"ethereum"`
}

type BinanceResponse struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

type KrakenResponse struct {
	Result struct {
		XXBTZUSD struct {
			C []string `json:"c"` // First element is the last trade price
		} `json:"XXBTZUSD"`
	} `json:"result"`
}

// FetchAPI pulls data from a given API URL and parses based on the source
func FetchAPI(name, url string) (*APIResponse, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching %s: %v", name, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse response based on the API source
	switch name {
	case "coingecko":
		return parseCoingeckoResponse(body)
	case "binance":
		return parseBinanceResponse(body)
	case "kraken":
		return parseKrakenResponse(body)
	default:
		return nil, fmt.Errorf("unknown API source: %s", name)
	}
}

func parseCoingeckoResponse(body []byte) (*APIResponse, error) {
	var response CoinGeckoResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &APIResponse{
		Timestamp: time.Now(),
		Source:    "coingecko",
		Symbol:    "BTC",
		Price:     response.Bitcoin.USD,
	}, nil
}

func parseBinanceResponse(body []byte) (*APIResponse, error) {
	var response BinanceResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	price, err := strconv.ParseFloat(response.Price, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing price: %v", err)
	}

	return &APIResponse{
		Timestamp: time.Now(),
		Source:    "binance",
		Symbol:    "BTC",
		Price:     price,
	}, nil
}

func parseKrakenResponse(body []byte) (*APIResponse, error) {
	var response KrakenResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	if len(response.Result.XXBTZUSD.C) == 0 {
		return nil, fmt.Errorf("no price data in Kraken response")
	}

	price, err := strconv.ParseFloat(response.Result.XXBTZUSD.C[0], 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing price: %v", err)
	}

	return &APIResponse{
		Timestamp: time.Now(),
		Source:    "kraken",
		Symbol:    "BTC",
		Price:     price,
	}, nil
}

// FetchFromMultipleAPIs fetches data from multiple sources concurrently
func FetchFromMultipleAPIs() []APIResponse {
	apiList := map[string]string{
		"coingecko": "https://api.coingecko.com/api/v3/simple/price?ids=bitcoin,ethereum&vs_currencies=usd",
		"binance":   "https://api.binance.com/api/v3/ticker/price?symbol=BTCUSDT",
		"kraken":    "https://api.kraken.com/0/public/Ticker?pair=XBTUSD",
	}

	results := make([]APIResponse, 0)
	ch := make(chan *APIResponse)

	// Fetch data concurrently
	for name, url := range apiList {
		go func(name, url string) {
			res, err := FetchAPI(name, url)
			if err == nil {
				ch <- res
			} else {
				fmt.Println("❌ API Fetch error:", err)
			}
		}(name, url)
	}

	// Collect results
	for i := 0; i < len(apiList); i++ {
		result := <-ch
		results = append(results, *result)
	}

	return results
}
