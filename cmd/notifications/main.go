package main

import (
	"context"
	"encoding/json"
	"fmt"
	"habr-app/internal/config"
	"log"

	"github.com/segmentio/kafka-go"
)

type RegistrationEvent struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

func main() {
	fmt.Println("Notification service is started...")
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{config.GetKafkaBroker()},
		Topic:       "user-registrations",
		StartOffset: kafka.FirstOffset,
		MinBytes:    1,
		MaxBytes:    10e6,
	})
	defer reader.Close()
	fmt.Println("Waiting messages...")
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			fmt.Printf("Error from reading message %v", err)
			continue
		}
		var event RegistrationEvent
		if err := json.Unmarshal(m.Value, &event); err != nil {
			log.Println("json unmarshal", err)
			continue
		}
		fmt.Printf("Sending verifying email to %s\n", event.Email)
		fmt.Printf("Verification link: %s/verify?token=%s\n", config.GetAuthHost(), event.Token)
	}
}
