.PHONY: build run clean test lint deps migrate-up migrate-down dev

# Application name
APP_NAME := auth
BUILD_DIR := bin
MAIN := cmd/main.go

# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod

build:
	$(GOBUILD) -o $(BUILD_DIR)/$(APP_NAME) $(MAIN)

run: build
	./$(BUILD_DIR)/$(APP_NAME)

clean:
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

test:
	$(GOTEST) -v ./...

deps:
	$(GOMOD) download
	$(GOMOD) tidy

lint:
	golangci-lint run

migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down

dev:
	$(GOBUILD) -o $(BUILD_DIR)/$(APP_NAME) $(MAIN)
	./$(BUILD_DIR)/$(APP_NAME)

help:
	@echo "Available commands:"
	@echo "  build      - Build the application"
	@echo "  run        - Build and run the application"
	@echo "  clean      - Clean build artifacts"
	@echo "  test       - Run tests"
	@echo "  deps       - Download and tidy dependencies"
	@echo "  lint       - Run linter"
	@echo "  migrate-up - Run database migrations"
	@echo "  migrate-down - Rollback database migrations"
	@echo "  dev        - Build and run for development"