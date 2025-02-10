package kafka

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	apis "crypto-monitor/internal/api"

	"github.com/IBM/sarama"
)

type Producer struct {
	Client sarama.SyncProducer
	Topic  string
}

// init and return producer
func NewProducer(brokers []string, topic string) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %v", err)
	}

	return &Producer{Client: producer, Topic: topic}, nil
}

// send message to kafka
func (p *Producer) SendMessage(apiData apis.APIResponse) error {
	msgJSON, err := json.Marshal(apiData)
	if err != nil {
		return fmt.Errorf("failed to encode message: %v", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: p.Topic,
		Value: sarama.StringEncoder(msgJSON),
	}

	// partition, offset, err := p.Client.SendMessage(msg)
	_, _, err = p.Client.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	return nil
}

// continuously fetch from apis and send to kafka
func (p *Producer) StartProducer() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		apiResults := apis.FetchFromMultipleAPIs()
		for _, result := range apiResults {
			if err := p.SendMessage(result); err != nil {
				log.Println(err)
			}
		}
	}
}
