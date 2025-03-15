PROJECT_DIR := $(shell pwd)
PROJECT_NAME := $(shell basename $(PROJECT_DIR))
SCRIPTS_DIR := $(PROJECT_DIR)/scripts

deps:
	@echo "Installing dependencies..."
	@pnpm install
	@echo "Done!"
.PHONY: deps

format:
	@echo "Formatting code..."
	@pnpm run format
	@echo "Done!"
.PHONY: format

gen\:validation:
	@echo "Generating validation schema..."
	@pnpm tsx $(SCRIPTS_DIR)/generateValidation.ts
	@echo "Done!"
.PHONY: gen\:validation

lint:
	@echo "Linting code..."
	@pnpm run lint
	@echo "Done!"
.PHONY: lint

start:
	@echo "Starting server..."
	@pnpm run start
	@echo "Done!"
.PHONY: start

test:
	@echo "Running tests..."
	@pnpm run test
	@echo "Done!"
.PHONY: test

tsc:
	@echo "Compiling typescript..."
	@pnpm run tsc
	@echo "Done!"
.PHONY: tsc
