package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/tv-anagha/ecommerce-backend/order-service/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func env(primary, fallback string) string {
	if v := os.Getenv(primary); v != "" {
		return v
	}
	return os.Getenv(fallback)
}

func Connect() {
	_ = godotenv.Load()
	_ = godotenv.Load("../.env")
	_ = godotenv.Load("../../.env")

	dbName := env("POSTGRES_DB", "DB_NAME")
	if dbName == "" {
		dbName = "order_db"
	}

	sslmode := env("POSTGRES_SSLMODE", "")
	if sslmode == "" {
		sslmode = "disable"
	}

	timezone := env("POSTGRES_TIMEZONE", "")
	if timezone == "" {
		timezone = "UTC"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s TimeZone=%s sslmode=%s",
		env("POSTGRES_HOST", "DB_HOST"),
		env("POSTGRES_USER", "DB_USER"),
		env("POSTGRES_PASSWORD", "DB_PASSWORD"),
		dbName,
		env("POSTGRES_PORT", "DB_PORT"),
		timezone,
		sslmode,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("order-service: database connect failed: %v", err)
	}

	if err := DB.AutoMigrate(&model.Order{}, &model.OrderItem{}); err != nil {
		log.Fatalf("order-service: migrate failed: %v", err)
	}

	log.Printf("order-service: database connected (dbname=%s)", dbName)
}
