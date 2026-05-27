# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Binary name and paths
BINARY_NAME=godaml
BINARY_DIR=bin
CMD_DIR=cmd
MAIN_PATH=./$(CMD_DIR)/go-daml/main.go

# Pinned tool versions
GOLANGCI_LINT_VERSION=v2.12.2
GOLANGCI_LINT=$(BINARY_DIR)/golangci-lint

# Smoke test paths
SMOKE_DIR=$(BINARY_DIR)/smoke
SMOKE_DARS=$(wildcard test-data/*.dar)
SMOKE_PKG=smoketest

# Build flags
LDFLAGS=-ldflags "-s -w"
BUILD_FLAGS=-trimpath

# Default target
.PHONY: all
all: clean build

# Build the CLI tool
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BINARY_DIR)
	$(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BINARY_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BINARY_DIR)/$(BINARY_NAME)"

# Build for multiple platforms
.PHONY: build-all
build-all: clean
	@echo "Building for multiple platforms..."
	@mkdir -p $(BINARY_DIR)
	
	@echo "Building for Linux amd64..."
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BINARY_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	
	@echo "Building for Linux arm64..."
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BINARY_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_PATH)
	
	@echo "Building for macOS amd64..."
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BINARY_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	
	@echo "Building for macOS arm64..."
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BINARY_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	
	@echo "Building for Windows amd64..."
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BINARY_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
	
	@echo "Multi-platform build complete"

# Install the CLI tool to GOPATH/bin
.PHONY: install
install:
	@echo "Installing $(BINARY_NAME)..."
	$(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(GOPATH)/bin/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Installed $(BINARY_NAME) to $(GOPATH)/bin/"

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Smoke: build everything, run CLI codegen end-to-end on every test-data/*.dar
# fixture, ensure the output compiles.
.PHONY: smoke
smoke: build
	@echo "Smoke: go build ./..."
	$(GOBUILD) ./...
	@echo "Smoke: end-to-end codegen on $(words $(SMOKE_DARS)) DAR fixtures..."
	@set -e; for dar in $(SMOKE_DARS); do \
		echo "  -> $$dar"; \
		rm -rf $(SMOKE_DIR); \
		mkdir -p $(SMOKE_DIR); \
		printf 'module $(SMOKE_PKG)\n\ngo 1.24.2\n\nreplace github.com/noders-team/go-daml => $(CURDIR)\n' > $(SMOKE_DIR)/go.mod; \
		./$(BINARY_DIR)/$(BINARY_NAME) --dar $$dar --output $(SMOKE_DIR) --go_package $(SMOKE_PKG); \
		(cd $(SMOKE_DIR) && $(GOMOD) tidy && $(GOBUILD) ./...); \
	done
	@echo "Smoke: OK ($(words $(SMOKE_DARS)) fixtures)"

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	@rm -rf $(BINARY_DIR)
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

# Download dependencies
.PHONY: deps
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Generate proto files (if needed)
.PHONY: proto
proto:
	@echo "Generating proto files..."
	@if [ -d "proto" ]; then \
		protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/*.proto; \
	fi

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

# Install pinned golangci-lint into ./bin
.PHONY: lint-install
lint-install:
	@if [ ! -x "$(GOLANGCI_LINT)" ] || ! $(GOLANGCI_LINT) --version 2>/dev/null | grep -q "$(GOLANGCI_LINT_VERSION:v%=%)"; then \
		echo "Installing golangci-lint $(GOLANGCI_LINT_VERSION)..."; \
		mkdir -p $(BINARY_DIR); \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/$(GOLANGCI_LINT_VERSION)/install.sh \
			| sh -s -- -b $(BINARY_DIR) $(GOLANGCI_LINT_VERSION); \
	fi

# Lint code with pinned version
.PHONY: lint
lint: lint-install
	@echo "Linting code..."
	$(GOLANGCI_LINT) run

# Vet code
.PHONY: vet
vet:
	@echo "Vetting code..."
	$(GOCMD) vet ./...

# Run the CLI tool with example parameters
.PHONY: run
run: build
	@echo "Running $(BINARY_NAME) with example parameters..."
	./$(BINARY_DIR)/$(BINARY_NAME) --help

# Development workflow: format, vet, test, build
.PHONY: dev
dev: fmt vet test build

# Release workflow: clean, deps, test, build-all
.PHONY: release
release: clean deps test build-all

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all           - Clean and build"
	@echo "  build         - Build the CLI tool for current platform"
	@echo "  build-all     - Build for multiple platforms"
	@echo "  install       - Install CLI tool to GOPATH/bin"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  clean         - Clean build artifacts"
	@echo "  deps          - Download and tidy dependencies"
	@echo "  proto         - Generate proto files"
	@echo "  fmt           - Format code"
	@echo "  lint          - Lint code (requires golangci-lint)"
	@echo "  vet           - Vet code"
	@echo "  run           - Build and run with help"
	@echo "  dev           - Development workflow (fmt, vet, test, build)"
	@echo "  release       - Release workflow (clean, deps, test, build-all)"
	@echo "  help          - Show this help"