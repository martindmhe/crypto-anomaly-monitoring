package kafka

import (
	apis "crypto-monitor/internal/api"
	detection "crypto-monitor/internal/detection"
	twilio "crypto-monitor/internal/alerts"
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
)


type Consumer struct {
	Client sarama.Consumer
	Topic  string
}

var windows = make(map[string]*detection.RollingWindow)

// get or create rolling window for each api
func GetOrCreateRollingWindow(apiName string) *detection.RollingWindow {
	if _, exists := windows[apiName]; !exists {
		windows[apiName] = detection.NewRollingWindow(50) // around 5 secs?
	}
	return windows[apiName]
}


func NewConsumer(brokers []string, topic string) (*Consumer, error) {
	consumer, err := sarama.NewConsumer(brokers, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %v", err)
	}

	return &Consumer{Client: consumer, Topic: topic}, nil
}

// listen for messages and process api responses
func (c *Consumer) StartConsumer() {
	partitions, err := c.Client.Partitions(c.Topic)
	if err != nil {
		log.Fatal("Failed to get partitions:", err)
	}

	for _, partition := range partitions {
		pc, err := c.Client.ConsumePartition(c.Topic, partition, sarama.OffsetNewest)
		if err != nil {
			log.Fatal("Failed to start consumer:", err)
		}

		go func(pc sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				// log.Println("received message test")
				var apiData apis.APIResponse
				if err := json.Unmarshal(msg.Value, &apiData); err != nil {
					log.Println("error with message: ", err)
					continue
				}

				source := apiData.Source
				window := GetOrCreateRollingWindow(source)
				
				price := apiData.Price
				window.AddPrice(price)

				zAnomaly, bAnomaly := window.CheckAnomalies(price)

				if zAnomaly {
					fmt.Printf("sending alert for z-score anomaly for %s api at price %v\n", source, price)
					twilio.SendTwilioAlert(source, price, "Z-Score Anomaly")
				}
				if bAnomaly {
					fmt.Printf("sending alert for bollinger anomaly for %s api at price %v\n", source, price)
					twilio.SendTwilioAlert(source, price, "Bollinger Bands")
				}

				// have to process differently since api schemas are different
				switch apiData.Source {
				case "coingecko":
					fmt.Printf("🌍 Coingecko Data:\nPrice: %v\nTime: %v\n", apiData.Price, apiData.Timestamp)
				case "binance":
					fmt.Printf("📈 Binance Data:\nPrice: %v\nTime: %v\n", apiData.Price, apiData.Timestamp)
				case "kraken":
					fmt.Printf("⚡ Kraken Data:\nPrice: %v\nTime: %v\n", apiData.Price, apiData.Timestamp)
				default:
					// fmt.Printf("other", apiData.Price)
				}
			}
		}(pc)
	}

	select {} // just so it doesn't exit
}
