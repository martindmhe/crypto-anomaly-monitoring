package exchange

import (
	"time"
	"math/rand"
)

type ExchangeClient interface {
    GetMarketData(symbol string) (*MarketData, error)
    Name() string
}

type MarketData struct {
    Exchange    string    `json:"exchange"`
    Symbol      string    `json:"symbol"`
    Price       float64   `json:"price"`
    Volume      float64   `json:"volume"`
    Timestamp   time.Time `json:"timestamp"`
}

type BinanceClient struct {
	// apiKey    string
	// apiSecret string
}

func (b *BinanceClient) Name() string {
	return "Binance"
}

func (b *BinanceClient) GetMarketData(symbol string) (*MarketData, error) {

	// call Binance API to get the market data
	// var apiURL string = "https://api.open-meteo.com/v1/forecast?latitude=43.6532&longitude=79.3832&current=temperature_2m,relative_humidity_2m"

	// // fmt.Println("fetched")

	// resp, err := http.Get(apiURL)

	// if err != nil {
	// 	return nil, err
	// }

	// defer resp.Body.Close()

    // This is a placeholder - you'd actually call Binance's API here
    return &MarketData{
        Exchange:  "Binance",
        Symbol:    symbol,
        Price:    rand.Float64() * 100000,
        Volume:   rand.Float64() * 1000,
        Timestamp: time.Now(),
    }, nil
}

