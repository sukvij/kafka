package database

import (
	"fmt"
	"vijju/logs"
	"vijju/user"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// InitDB initializes the MySQL database connection and migrates schemas
func InitDB() (*gorm.DB, error) {
	dsn := "root:root@tcp(mysql:3306)/users?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL: %v", err)
	}

	// Migrate User and UserLog models
	if err := db.AutoMigrate(&user.User{}, &logs.UserLog{}); err != nil {
		return nil, fmt.Errorf("failed to migrate schema: %v", err)
	}

	return db, nil
}
