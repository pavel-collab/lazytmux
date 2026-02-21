# LazyTmux Makefile

BINARY_NAME=lazytmux
BUILD_DIR=bin
CMD_DIR=./cmd/lazytmux
COVERAGE_FILE=coverage.out

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
GOVET=$(GOCMD) vet

# Build flags
LDFLAGS=-ldflags "-s -w"

.PHONY: all build clean test coverage install uninstall fmt vet tidy help

## Build commands

all: clean build ## Clean and build

build: ## Build the binary
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)
	@echo "Binary built: $(BUILD_DIR)/$(BINARY_NAME)"

build-dev: ## Build without optimization (for debugging)
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)

## Install/Uninstall

install: build ## Install to $GOPATH/bin
	cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/$(BINARY_NAME)
	@echo "Installed to $(GOPATH)/bin/$(BINARY_NAME)"

uninstall: ## Remove from $GOPATH/bin
	rm -f $(GOPATH)/bin/$(BINARY_NAME)
	@echo "Uninstalled from $(GOPATH)/bin/$(BINARY_NAME)"

## Testing

test: ## Run all tests
	$(GOTEST) -v ./...

test-short: ## Run tests in short mode (skip integration tests)
	$(GOTEST) -v -short ./...

coverage: ## Run tests with coverage
	$(GOTEST) -coverprofile=$(COVERAGE_FILE) ./...
	$(GOCMD) tool cover -func=$(COVERAGE_FILE)

coverage-html: coverage ## Generate HTML coverage report
	$(GOCMD) tool cover -html=$(COVERAGE_FILE) -o coverage.html
	@echo "Coverage report: coverage.html"

## Code quality

fmt: ## Format code
	$(GOCMD) fmt ./...

vet: ## Run go vet
	$(GOVET) ./...

lint: fmt vet ## Run all linters

## Dependencies

tidy: ## Tidy and verify go modules
	$(GOMOD) tidy
	$(GOMOD) verify

deps: ## Download dependencies
	$(GOMOD) download

## Cleanup

clean: ## Remove build artifacts
	rm -rf $(BUILD_DIR)
	rm -f $(COVERAGE_FILE)
	rm -f coverage.html

## Help

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

# Default target
.DEFAULT_GOAL := help
