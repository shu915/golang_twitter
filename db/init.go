// db/db.go
package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

var DB *pgx.Conn

func Init() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: .env file not found, using environment variables: %v", err)
	}

	DB_USER := os.Getenv("DB_USER")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_HOST := os.Getenv("DB_HOST")
	if DB_HOST == "" {
		DB_HOST = "localhost"
	}
	DB_NAME := os.Getenv("DB_NAME")
	DB_URL := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", DB_USER, DB_PASSWORD, DB_HOST, DB_NAME)

	DB, err = pgx.Connect(context.Background(), DB_URL)
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}

	// 接続確認
	if err := DB.Ping(context.Background()); err != nil {
		log.Fatalf("DB ping failed: %v", err)
	}
}
