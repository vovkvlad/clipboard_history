# Binary name
BINARY_NAME=clipboard_history
BINARY_DIR=.bin

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BINARY_DIR)
	@go build -o $(BINARY_DIR)/$(BINARY_NAME) .

# Build for different platforms
build-linux:
	@echo "Building for Linux..."
	@mkdir -p $(BINARY_DIR)
	@GOOS=linux GOARCH=amd64 go build -o $(BINARY_DIR)/$(BINARY_NAME)-linux .

build-windows:
	@echo "Building for Windows..."
	@mkdir -p $(BINARY_DIR)
	@GOOS=windows GOARCH=amd64 go build -o $(BINARY_DIR)/$(BINARY_NAME).exe .

build-darwin:
	@echo "Building for macOS..."
	@mkdir -p $(BINARY_DIR)
	@GOOS=darwin GOARCH=amd64 go build -o $(BINARY_DIR)/$(BINARY_NAME)-darwin .

# Build for all platforms
build-all: build-linux build-windows build-darwin

# Run the application
run:
	@echo "Running $(BINARY_NAME)..."
	@go run .

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BINARY_DIR)

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

# Run tests
test:
	@echo "Running tests..."
	@go test ./...

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	@golangci-lint run

# Show help
help:
	@echo "Available commands:"
	@echo "  build       - Build the application"
	@echo "  build-linux - Build for Linux"
	@echo "  build-windows - Build for Windows"
	@echo "  build-darwin - Build for macOS"
	@echo "  build-all   - Build for all platforms"
	@echo "  run         - Run the application"
	@echo "  clean       - Clean build artifacts"
	@echo "  deps        - Install dependencies"
	@echo "  test        - Run tests"
	@echo "  fmt         - Format code"
	@echo "  lint        - Lint code"
	@echo "  help        - Show this help"

.PHONY: build build-linux build-windows build-darwin build-all run clean deps test fmt lint help
