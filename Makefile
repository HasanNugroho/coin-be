.PHONY: help build run dev swagger swagger-gen clean install-tools seed docker-build docker-up docker-down docker-logs docker-seed

help:
	@echo "Available commands:"
	@echo ""
	@echo "Development:"
	@echo "  make install-tools    - Install Air and Swagger tools"
	@echo "  make build            - Build the application"
	@echo "  make run              - Run the application"
	@echo "  make dev              - Run with hot reload (Air)"
	@echo "  make swagger-gen      - Generate Swagger documentation"
	@echo "  make swagger          - Generate and open Swagger docs"
	@echo "  make seed             - Seed database with default data"
	@echo ""
	@echo "Bot:"
	@echo "  make build-bot        - Build the Telegram bot"
	@echo "  make run-bot          - Run the Telegram bot"
	@echo "  make bot              - Build and run the Telegram bot"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-build     - Build Docker image"
	@echo "  make docker-up        - Start Docker containers"
	@echo "  make docker-down      - Stop Docker containers"
	@echo "  make docker-logs      - View Docker logs"
	@echo "  make docker-seed      - Seed database in Docker"
	@echo "  make docker-clean     - Remove Docker containers and volumes"
	@echo ""
	@echo "Utilities:"
	@echo "  make clean            - Clean build artifacts"
	@echo "  make test             - Run tests"
	@echo "  make fmt              - Format code"
	@echo "  make lint             - Run linter"

install-tools:
	@echo "Installing Air and Swagger tools..."
	go install github.com/air-verse/air@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	@echo "Tools installed successfully"

build:
	@echo "Building application..."
	go build -o bin/main ./cmd/api
	@echo "Build complete: bin/main"

run: build
	@echo "Running application..."
	./bin/main

build-bot:
	@echo "Building Telegram bot..."
	go build -o bin/bot ./cmd/bot
	@echo "Build complete: bin/bot"

run-bot: build-bot
	@echo "Running Telegram bot..."
	./bin/bot

bot: run-bot

dev:
	@echo "Starting development server with hot reload..."
	air -c .air.toml

swagger-gen:
	@echo "Generating Swagger documentation..."
	swag init -g cmd/api/main.go --output docs
	@echo "Swagger documentation generated"

swagger: swagger-gen
	@echo "Swagger docs available at http://localhost:8080/swagger/index.html"

clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -rf tmp/
	rm -f build-errors.log
	@echo "Clean complete"

test:
	@echo "Running tests..."
	go test -v ./...

fmt:
	@echo "Formatting code..."
	go fmt ./...
	@echo "Format complete"

lint:
	@echo "Running linter..."
	golangci-lint run ./...

deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy
	@echo "Dependencies downloaded"

seed:
	@echo "Seeding database with default data..."
	go run cmd/seeder/main.go
	@echo "Database seeding complete"

docker-build:
	@echo "Building Docker image..."
	docker-compose build
	@echo "Docker image built successfully"

docker-up:
	@echo "Starting Docker containers..."
	docker-compose up -d
	@echo "Waiting for services to be healthy..."
	@sleep 5
	@echo "Docker containers started successfully"
	@echo "Application: http://localhost:8080"
	@echo "Swagger: http://localhost:8080/swagger/index.html"
	@echo "MongoDB: localhost:27017"
	@echo "Redis: localhost:6379"

docker-down:
	@echo "Stopping Docker containers..."
	docker-compose down
	@echo "Docker containers stopped"

docker-logs:
	@echo "Showing Docker logs..."
	docker-compose logs -f

docker-seed:
	@echo "Seeding database in Docker..."
	docker-compose exec app ./seeder
	@echo "Database seeding complete"

docker-clean:
	@echo "Removing Docker containers and volumes..."
	docker-compose down -v
	@echo "Docker cleanup complete"

.DEFAULT_GOAL := help
