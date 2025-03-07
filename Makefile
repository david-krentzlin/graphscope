.PHONY: build test clean run

# Variables
BINARY_NAME=graphscope
BUILD_DIR=bin

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/graphscope

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)

# Run the application
run: build
	@./$(BUILD_DIR)/$(BINARY_NAME)

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod tidy

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
