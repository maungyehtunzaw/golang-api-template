package main

import (
	"context"
	"log"
	"time"

	"golang-api-template/internal/config"
	"golang-api-template/internal/migrations"
	"golang-api-template/internal/router"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Load main config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Setup DB
	db, err := config.SetupDatabase(cfg)
	if err != nil {
		log.Fatalf("DB setup error: %v", err)
	}
	if err := migrations.AutoMigrateDatabase(db); err != nil {
		log.Fatalf("Error migrating database: %v", err)
	}

	// emailConfig := config.GetEmailConfig()
	// emailService := service.NewEmailService(emailConfig)

	// // Example: Send a test email upon startup
	// err100 := emailService.SendEmail("test@example.com", "Test Email", "Hello, this is a test email from our application.")
	// if err100 != nil {
	// 	log.Fatalf("Failed to send test email: %v", err)
	// }

	// Setup Redis
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	redisClient, err := config.SetupRedis(ctx, cfg.Redis)
	if err != nil {
		log.Fatalf("Redis setup error: %v", err)
	}
	defer redisClient.Close()

	// Router
	r := router.SetupRouter(db, redisClient, cfg)

	log.Printf("Starting server on port %s...", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
