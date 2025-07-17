include .env
export

DB_URL=postgres://${DB_USER}:${DB_PASSWORD}@localhost:5432/${DB_NAME}?sslmode=disable
MIGRATE=migrate -path db/migration -database $(DB_URL)

.PHONY: migrate-up migrate-down migrate-create

migrate-up:
	$(MIGRATE) up

migrate-down:
	$(MIGRATE) down

migrate-create:
	migrate create -ext sql -dir db/migration -seq $(name)