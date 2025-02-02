package main

import (
	_ "fiber-app/docs" // swagger docs
	"fiber-app/pkg/cache"
	"fiber-app/pkg/cron"
	"fiber-app/pkg/database"
	"fiber-app/pkg/handlers"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
)

// @title Fiber Message API
// @version 1.0
// @description Bu API, mesaj gönderme ve yönetme işlemleri için kullanılır.
// @termsOfService http://swagger.io/terms/

// @contact.name API Destek
// @contact.email destek@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /api
func main() {
	app := fiber.New()

	app.Use(cors.New())

	if err := database.Connect(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Redis bağlantısını başlat
	if err := cache.Connect(); err != nil {
		log.Printf("Warning: Failed to initialize Redis: %v", err)
	}

	// Swagger endpoint'i
	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL:         "/swagger/doc.json",
		DeepLinking: true,
	}))

	api := app.Group("/api")
	api.Post("/messages", handlers.CreateMessage)
	api.Get("/messages", handlers.GetMessages)
	api.Post("/cron/start", handlers.StartCronJob)
	api.Post("/cron/stop", handlers.StopCronJob)
	api.Get("/cron/status", handlers.GetCronStatus)
	api.Get("/cron/logs", handlers.GetCronLogs)

	// Start cron job by default
	if err := cron.StartCron(); err != nil {
		log.Printf("Warning: Failed to start cron job: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Fatal(app.Listen(":" + port))
}
