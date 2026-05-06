package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"habr-app/internal/config"

	"github.com/segmentio/kafka-go"
)

type RegistrationEvent struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

func SendRegistrationEvent(email, token string) error {
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(config.GetKafkaBroker()),
		Topic:                  "user-registrations",
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
	}
	defer writer.Close()

	event := RegistrationEvent{
		Email: email,
		Token: token,
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("json marshal error: %w", err)
	}

	err = writer.WriteMessages(context.Background(),
		kafka.Message{
			Value: payload,
		},
	)
	if err != nil {
		return fmt.Errorf("kafka write error: %w", err)
	}

	fmt.Printf("Message sent to Kafka: email=%s\n", email)
	return nil
}