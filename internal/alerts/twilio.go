package alerts

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

func SendTwilioAlert(source string, price float64, anomalyType string) {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
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

	fmt.Printf("Alert sent for %s api at price %v\n", source, price)
}
