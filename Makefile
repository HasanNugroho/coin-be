.PHONY: help build run dev swagger swagger-gen clean install-tools

help:
	@echo "Available commands:"
	@echo "  make install-tools    - Install Air and Swagger tools"
	@echo "  make build            - Build the application"
	@echo "  make run              - Run the application"
	@echo "  make dev              - Run with hot reload (Air)"
	@echo "  make swagger-gen      - Generate Swagger documentation"
	@echo "  make swagger          - Generate and open Swagger docs"
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

.DEFAULT_GOAL := help
