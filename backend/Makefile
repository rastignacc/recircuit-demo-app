.PHONY: run build test migrate-up migrate-down

DATABASE_URL ?= postgres://marketplace:marketplace@localhost:5432/marketplace?sslmode=disable

run:
	go run ./cmd/server

build:
	go build -o bin/server ./cmd/server

test:
	go test ./...

migrate-up:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate \
		-path ./migrations -database "$(DATABASE_URL)" up

migrate-down:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate \
		-path ./migrations -database "$(DATABASE_URL)" down 1
