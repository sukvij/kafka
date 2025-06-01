package main

import (
	"log"
	"vijju/database"
	"vijju/logs"
	"vijju/user"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Start Kafka log consumer
	go logs.StartLogConsumer(db, []string{"user-logs"})

	// Initialize Gin router
	r := gin.Default()

	// Initialize user controller
	userController := user.NewController(db)

	// Define routes
	r.POST("/users", userController.CreateUser)
	r.GET("/users/:id", userController.GetUser)
	r.PUT("/users/:id", userController.UpdateUser)
	r.DELETE("/users/:id", userController.DeleteUser)

	// Start server
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
