package logs

import (
	"time"
	"vijju/kafka"

	"gorm.io/gorm"
)

// UserLog model for MySQL user_logs table
type UserLog struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `json:"user_id"`
	Action    string    `json:"action"` // created, updated, deleted
	Timestamp time.Time `json:"timestamp"`
}

// StartLogConsumer starts Kafka consumers for specified topics and inserts logs into the database
func StartLogConsumer(db *gorm.DB, topics []string) {
	for _, topic := range topics {
		go kafka.Reader(topic, func(event kafka.UserEvent) error {
			logEntry := UserLog{
				UserID:    event.UserID,
				Action:    event.Action,
				Timestamp: event.Timestamp,
			}
			return db.Create(&logEntry).Error
		})
	}
}
