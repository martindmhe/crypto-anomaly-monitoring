package alerts

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"

	"time"

	"github.com/go-redis/redis"
)

var redisClient *redis.Client

func InitRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	for {
		_, err := redisClient.Ping().Result()
		if err == nil {
			log.Println("Redis connected successfully")
			break
		}
		log.Printf("Waiting for Redis... %v", err)
		time.Sleep(1 * time.Second)
	}
}

func CanSendAlert(source string) bool {
	key := fmt.Sprintf("last_alert:%s", source)
	lastAlertTime, err := redisClient.Get(key).Result()
	if err == redis.Nil {
		// since this means key doesn't exist
		return true
	} else if err != nil {
		log.Printf("error getting last time: %v", err)
		return false
	}

	lastTime, err := time.Parse(time.RFC3339, lastAlertTime)
	if err != nil {
		log.Println("parsing error:", err)
		return false
	}

	timeSinceLastAlert := time.Since(lastTime)
	if timeSinceLastAlert < 30*time.Second {
		return false
	}

	return true

}

func SendTwilioAlert(source string, price float64, anomalyType string) {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	if !CanSendAlert(source) {
		return
	}

	client := twilio.NewRestClient()

	toPhone := os.Getenv("TWILIO_TO_PHONE")
	fromPhone := os.Getenv("TWILIO_FROM_PHONE")
	body := fmt.Sprintf("🚨 %s API Alert: %s anomaly detected at price %.2f", source, anomalyType, price)

	params := &openapi.CreateMessageParams{
		To:   &toPhone,
		From: &fromPhone,
		Body: &body,
	}

	_, err := client.Api.CreateMessage(params)
	if err != nil {
		log.Printf("Error sending message: %v", err)
		return
	}

	key := fmt.Sprintf("last_alert:%s", source)
	redisClient.Set(key, time.Now().Format(time.RFC3339), 30*time.Second)

	fmt.Printf("Alert sent for %s api at price %v\n", source, price)
}
