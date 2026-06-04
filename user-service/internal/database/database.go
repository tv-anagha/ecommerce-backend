package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
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

	sslmode := env("POSTGRES_SSLMODE", "")
	if sslmode == "" {
		sslmode = "disable"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s TimeZone=%s sslmode=%s",
		env("POSTGRES_HOST", "DB_HOST"),
		env("POSTGRES_USER", "DB_USER"),
		env("POSTGRES_PASSWORD", "DB_PASSWORD"),
		env("POSTGRES_DB", "DB_NAME"),
		env("POSTGRES_PORT", "DB_PORT"),
		env("POSTGRES_TIMEZONE", "UTC"),
		sslmode,
	)

	var err error

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connected successfully")
}