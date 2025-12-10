.PHONY: build run test test-coverage lint clean migrate

BINARY_NAME=vacaytracker-api
BUILD_DIR=./bin
GOFLAGS=-ldflags="-w -s"

## build: Build the application binary
build:
	CGO_ENABLED=0 go build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/server

## run: Run the application
run:
	go run ./cmd/server/main.go

## test: Run all tests
test:
	go test -v ./...

## test-coverage: Run tests with coverage report
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

## lint: Run linter
lint:
	golangci-lint run ./...

## clean: Clean build artifacts
clean:
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

## migrate: Run database migrations
migrate:
	go run ./cmd/server/main.go -migrate

## deps: Download dependencies
deps:
	go mod download
	go mod tidy

## help: Show this help message
help:
	@echo "Available targets:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'
