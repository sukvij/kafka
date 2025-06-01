package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

// UserEvent for Kafka messages
type UserEvent struct {
	UserID    uint      `json:"user_id"`
	Action    string    `json:"action"`
	Timestamp time.Time `json:"timestamp"`
}

// Writer writes messages to a Kafka topic
func Writer(userID uint, action string) error {
	writer := &kafka.Writer{
		Addr:     kafka.TCP("kafka:9092"),
		Topic:    "user-logs",
		Balancer: &kafka.Hash{},
	}
	defer writer.Close()

	event := UserEvent{
		UserID:    userID,
		Action:    action,
		Timestamp: time.Now(),
	}
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %v", err)
	}

	err = writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(fmt.Sprintf("%d", userID)),
		Value: eventBytes,
	})
	if err != nil {
		return fmt.Errorf("failed to write to Kafka: %v", err)
	}

	return nil
}

// Reader reads messages from a Kafka topic and processes them
func Reader(topic string, processFunc func(UserEvent) error) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"kafka:9092"},
		Topic:    topic,
		GroupID:  "user-log-group",
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
		MaxWait:  100 * time.Millisecond,
	})
	defer reader.Close()

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Consumer error for topic %s: %v", topic, err)
			continue
		}

		var event UserEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("Failed to unmarshal event from topic %s: %v", topic, err)
			continue
		}

		if err := processFunc(event); err != nil {
			log.Printf("Failed to process event from topic %s: %v", topic, err)
		} else {
			log.Printf("Processed event from topic %s: UserID=%d, Action=%s", topic, event.UserID, event.Action)
		}
	}
}
