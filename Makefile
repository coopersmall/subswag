PROJECT_DIR := $(shell pwd)

BIN_DIR := $(PROJECT_DIR)/bin
CMD_DIR := $(PROJECT_DIR)/cmd

PROJECT_NAME := $(shell basename $(PROJECT_DIR))
PROJECT_MAIN := $(CMD_DIR)/$(PROJECT_NAME)/api/main.go
SCRIPTS_DIR := $(PROJECT_DIR)/scripts

GOBIN := $(shell go env GOPATH)/bin
export PATH := $(GOBIN):$(PATH)

build:
	@echo "Building $(PROJECT_NAME)"
	@go build -o bin/$(PROJECT_NAME) $(PROJECT_MAIN)
	@echo "Done!"
.PHONY: build

clean:
	@echo "Cleaning $(PROJECT_NAME)"
	@rm -rf $(BIN_DIR)
	@echo "Done!"
.PHONY: clean

deps: gen\:deps
	@echo "Installing Dependencies"
	@chmod +x $(SCRIPTS_DIR)/db.sh
	@chmod +x $(SCRIPTS_DIR)/test.sh
	@go install github.com/air-verse/air@latest
	@go mod tidy
	@echo "Done!"
.PHONY: deps

docker-compose:
	@echo "Starting Docker Compose"
	@docker-compose up -d
	@echo "Done!"

format:
	@echo "Formatting $(PROJECT_NAME)"
	@go fmt ./...
	@echo "Done!"
.PHONY: format

gen\:data:
	@echo "Setting up $(PROJECT_NAME)"
	@go run $(SCRIPTS_DIR)/gen/local/data/main.go
	@echo "Done!"
.PHONY: gen\:data

gen\:deps:
	@echo "Generating Dependencies"
	@go run $(SCRIPTS_DIR)/gen/local/deps/main.go
	@echo "Done!"
.PHONY: gen\:deps

gen\:sql:
	@echo "Generating SQL"
	@sqlc generate
	@echo "Done!"
.PHONY: gen\:sql

gen\:ts:
	@echo "Generating TypeScript Types"
	@go run $(SCRIPTS_DIR)/gen/local/typescript/main.go
	@echo "Done!"
.PHONY: gen\:ts

lint:
	@echo "Linting $(PROJECT_NAME)"
	@go vet ./...
	@echo "Done!"
.PHONY: lint

psql:
	@echo "Connecting to PostgreSQL"
	@chmod +x $(SCRIPTS_DIR)/db.sh
	@$(SCRIPTS_DIR)/db.sh
.PHONY: psql

start: start\:db
	@echo "Starting $(PROJECT_NAME)"
	@echo
	@go run $(PROJECT_MAIN)
	@echo "Done!"
.PHONY: start

start\:w:
	@echo "Starting $(PROJECT_NAME)"
	@echo
	@air
	@echo "Done!"

start\:db:
	@echo "Starting PostgreSQL"
	@docker-compose up -d db
	@echo "Done!"
.PHONY: start\:db

start\:cache:
	@echo "Starting Redis"
	@docker-compose up -d cache
	@echo "Done!"
.PHONY: start\:cache

start\:web: gen\:web
	@echo "Generating HTML"
	@make gen\:web
	@echo "Starting Web on port 9000"
	@go run $(PROJECT_MAIN)
	@echo "Done!"
.PHONY: start\:web

test:
	@echo "Running all tests initially"
	@gotest -cover ./...
	@echo "Done"
.PHONY: test

test\:w:
	@echo "Running all tests initially"
	@$(SCRIPTS_DIR)/test.sh
	@echo "Done"
.PHONY: test\:w

test\:v:
	@echo "Running Tests"
	@go test -cover -v ./...
	@echo "Done"
.PHONY: test\:verbose
