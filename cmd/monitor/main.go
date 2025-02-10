package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"crypto-monitor/internal/alerts"
	"crypto-monitor/internal/kafka"
)

func main() {
	brokers := []string{"localhost:9092"}
	topic := "crypto-prices"

	producer, err := kafka.NewProducer(brokers, topic)
	if err != nil {
		log.Fatal(err)
	}
	defer producer.Client.Close()

	consumer, err := kafka.NewConsumer(brokers, topic)
	if err != nil {
		log.Fatal(err)
	}
	defer consumer.Client.Close()

	// start redis server for rate limiting (dont blow up twilio api)
	alerts.InitRedis()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		producer.StartProducer()
	}()

	go func() {
		defer wg.Done()
		consumer.StartConsumer()
	}()

	// log.Println("test")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	log.Println("stopping ")
	wg.Wait()
}
