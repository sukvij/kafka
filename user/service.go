package user

import (
	"log"
	"strconv"
	"vijju/kafka"

	"gorm.io/gorm"
)

// Service handles user business logic
type Service struct {
	DB *gorm.DB
}

// NewService creates a new user service
func NewService(db *gorm.DB) *Service {
	return &Service{DB: db}
}

// CreateUser creates a user and produces a Kafka log
func (s *Service) CreateUser(user *User) error {
	if err := s.DB.Create(user).Error; err != nil {
		return err
	}

	if err := kafka.Writer(user.ID, "created"); err != nil {
		log.Printf("Failed to write Kafka log: %v", err)
	}

	return nil
}

// GetUser retrieves a user by ID
func (s *Service) GetUser(id string) (*User, error) {
	var user User
	if err := s.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates a user and produces a Kafka log
func (s *Service) UpdateUser(id string, input *User) (*User, error) {
	var user User
	if err := s.DB.First(&user, id).Error; err != nil {
		return nil, err
	}

	user.Email = input.Email
	user.Name = input.Name
	if err := s.DB.Save(&user).Error; err != nil {
		return nil, err
	}

	if err := kafka.Writer(user.ID, "updated"); err != nil {
		log.Printf("Failed to write Kafka log: %v", err)
	}

	return &user, nil
}

// DeleteUser deletes a user and produces a Kafka log
func (s *Service) DeleteUser(id string) error {
	var user User
	if err := s.DB.First(&user, id).Error; err != nil {
		return err
	}

	if err := s.DB.Delete(&user).Error; err != nil {
		return err
	}

	userID, _ := strconv.ParseUint(id, 10, 32)
	if err := kafka.Writer(uint(userID), "deleted"); err != nil {
		log.Printf("Failed to write Kafka log: %v", err)
	}

	return nil
}
