// db/db.go
package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)


func Init() (*pgxpool.Pool, error) {
	_ = godotenv.Load()

	DB_USER := os.Getenv("DB_USER")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_HOST := os.Getenv("DB_HOST")
	DB_NAME := os.Getenv("DB_NAME")
	DB_URL := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", DB_USER, DB_PASSWORD, DB_HOST, DB_NAME)

	DB, err := pgxpool.New(context.Background(), DB_URL)
	if err != nil {
		return nil, err
	}

	// 接続確認
	if err := DB.Ping(context.Background()); err != nil {
		return nil, err
	}

	return DB, nil
}
