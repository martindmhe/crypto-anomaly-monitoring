package kafka

import (
	apis "crypto-monitor/internal/api"
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

// Consumer struct holds Kafka consumer instance
type Consumer struct {
	Client sarama.Consumer
	Topic  string
}

// NewConsumer initializes a Kafka consumer
func NewConsumer(brokers []string, topic string) (*Consumer, error) {
	consumer, err := sarama.NewConsumer(brokers, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %v", err)
	}

	return &Consumer{Client: consumer, Topic: topic}, nil
}

// StartConsumer listens for messages & processes API responses
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
				var apiData apis.APIResponse
				if err := json.Unmarshal(msg.Value, &apiData); err != nil {
					log.Println("❌ Failed to decode message:", err)
					continue
				}

				// Process API data differently based on source
				switch apiData.Source {
				case "coingecko":
					fmt.Printf("🌍 Coingecko Data:\nPrice: %v\nTime: %v\n", apiData.Price, apiData.Timestamp)
				case "binance":
					fmt.Printf("📈 Binance Data:\nPrice: %v\nTime: %v\n", apiData.Price, apiData.Timestamp)
				case "kraken":
					fmt.Printf("⚡ Kraken Data:\nPrice: %v\nTime: %v\n", apiData.Price, apiData.Timestamp)
				default:
					// fmt.Printf("🤖 Unknown API Data:", apiData.Price)
				}
			}
		}(pc)
	}

	select {} // Keep running
}
