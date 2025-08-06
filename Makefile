# Makefile for Hippodamus

# Variables
BINARY_NAME=hippodamus
BINARY_WINDOWS=$(BINARY_NAME).exe
BINARY_LINUX=$(BINARY_NAME)
BINARY_DARWIN=$(BINARY_NAME)

# Default target
.PHONY: all
all: test build

# Build the application
.PHONY: build
build:
	go build -o $(BINARY_WINDOWS) ./cmd/hippodamus

# Build for multiple platforms
.PHONY: build-all
build-all: build-windows build-linux build-darwin

.PHONY: build-windows
build-windows:
	GOOS=windows GOARCH=amd64 go build -o dist/$(BINARY_WINDOWS) ./cmd/hippodamus

.PHONY: build-linux
build-linux:
	GOOS=linux GOARCH=amd64 go build -o dist/$(BINARY_LINUX) ./cmd/hippodamus

.PHONY: build-darwin
build-darwin:
	GOOS=darwin GOARCH=amd64 go build -o dist/$(BINARY_DARWIN) ./cmd/hippodamus

# Run tests
.PHONY: test
test:
	go test -v ./...

# Clean build artifacts
.PHONY: clean
clean:
	go clean
	rm -f $(BINARY_WINDOWS)
	rm -rf dist/

# Install dependencies
.PHONY: deps
deps:
	go mod download
	go mod tidy

# Run examples
.PHONY: examples
examples: build
	@echo "Running infrastructure example..."
	./$(BINARY_WINDOWS) -input examples/infrastructure.yaml -templates examples/templates -output examples/infrastructure.drawio -verbose
	@echo "Running microservices example..."
	./$(BINARY_WINDOWS) -input examples/microservices.yaml -templates examples/templates -output examples/microservices.drawio -verbose
	@echo "Running simple example..."
	./$(BINARY_WINDOWS) -input examples/simple.yaml -output examples/simple.drawio

# Validate examples
.PHONY: validate
validate: build
	@echo "Validating infrastructure example..."
	./$(BINARY_WINDOWS) -validate -input examples/infrastructure.yaml -templates examples/templates
	@echo "Validating microservices example..."
	./$(BINARY_WINDOWS) -validate -input examples/microservices.yaml -templates examples/templates
	@echo "Validating simple example..."
	./$(BINARY_WINDOWS) -validate -input examples/simple.yaml

# Format code
.PHONY: fmt
fmt:
	go fmt ./...

# Lint code
.PHONY: lint
lint:
	golangci-lint run

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build       - Build the application"
	@echo "  build-all   - Build for all platforms"
	@echo "  test        - Run tests"
	@echo "  clean       - Clean build artifacts"
	@echo "  deps        - Install dependencies"
	@echo "  examples    - Run all examples"
	@echo "  validate    - Validate all examples"
	@echo "  fmt         - Format code"
	@echo "  lint        - Lint code"
	@echo "  help        - Show this help"
