APP_NAME = scality-cosi-driver
BIN_DIR = ./bin

.PHONY: all build test clean

all: test build

build:
	@echo "Building $(APP_NAME)..."
	go build -o $(BIN_DIR)/$(APP_NAME) ./cmd/$(APP_NAME)

test:
	@echo "Running tests..."
	go test ./... -coverprofile=coverage.txt -covermode=atomic

clean:
	@echo "Cleaning up..."
	rm -rf $(BIN_DIR)
