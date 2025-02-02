package database

import (
	"fiber-app/pkg/models"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() error {
	var err error

	// Get environment variables with default values
	dbHost := getEnvWithDefault("DB_HOST", "mysql")
	dbPort := getEnvWithDefault("DB_PORT", "3306")
	dbUser := getEnvWithDefault("DB_USER", "user")
	dbPass := getEnvWithDefault("DB_PASSWORD", "password")
	dbName := getEnvWithDefault("DB_NAME", "messages_db")

	// Log environment variables (without password)
	log.Printf("Database Configuration - Host: %s, Port: %s, User: %s, Database: %s\n",
		dbHost, dbPort, dbUser, dbName)

	// MySQL DSN format: username:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser,
		dbPass,
		dbHost,
		dbPort,
		dbName)

	log.Printf("Attempting to connect to database with DSN: %s\n", fmt.Sprintf("%s:****@tcp(%s:%s)/%s", dbUser, dbHost, dbPort, dbName))

	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("Failed to connect to database (attempt %d/%d): %v\n", i+1, maxRetries, err)
		if i < maxRetries-1 {
			time.Sleep(5 * time.Second)
		}
	}

	if err != nil {
		log.Printf("Failed to connect to database after %d attempts: %v\n", maxRetries, err)
		return err
	}

	// Drop existing tables
	if err := DB.Migrator().DropTable(&models.Message{}, &models.CronLog{}); err != nil {
		log.Printf("Failed to drop tables: %v\n", err)
		return err
	}
	log.Println("Existing tables dropped successfully")

	// Create tables
	if err := DB.AutoMigrate(&models.Message{}, &models.CronLog{}); err != nil {
		log.Printf("Failed to create tables: %v\n", err)
		return err
	}
	log.Println("Tables created successfully")

	// Insert default messages
	defaultMessages := []models.Message{
		{Content: "Hello! How can I help you?", Phone: "+905551234567", Status: false},
		{Content: "Good day, your order is being prepared.", Phone: "+905551234568", Status: false},
		{Content: "Your order has been shipped, it will arrive soon.", Phone: "+905551234569", Status: false},
		{Content: "Would you like to be informed about our campaigns?", Phone: "+905551234570", Status: false},
		{Content: "Update your profile for exclusive discount opportunities.", Phone: "+905551234571", Status: false},
	}

	if err := DB.Create(&defaultMessages).Error; err != nil {
		log.Printf("Failed to insert default messages: %v\n", err)
		return err
	}
	log.Printf("Successfully inserted %d default messages\n", len(defaultMessages))

	log.Println("Database connection established and initialized successfully")
	return nil
}

// getEnvWithDefault returns environment variable value or default if not set
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
