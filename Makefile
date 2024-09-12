APP_NAME = scality-cosi-driver

.PHONY: all
all: test build

.PHONY: build
build:
	@echo "Building $(APP_NAME)..."
	GOOS=linux GOARCH=amd64 go build -o ./bin/$(APP_NAME) ./cmd/$(APP_NAME)

.PHONY: test
test:
	@echo "Running tests..."
	go test ./... -coverprofile=coverage.txt -covermode=atomic

.PHONY: clean
clean:
	@echo "Cleaning up..."
	rm -rf ./bin
