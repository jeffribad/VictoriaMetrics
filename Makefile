# Makefile for VictoriaMetrics
# Provides common build, test, and deployment targets

APP_NAME := victoria-metrics
VERSION ?= $(shell git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.1")
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')

GO := go
GOFLAGS := -trimpath
LD_FLAGS := -s -w \
	-X main.Version=$(VERSION) \
	-X main.GitCommit=$(GIT_COMMIT) \
	-X main.BuildTime=$(BUILD_TIME)

BIN_DIR := bin
CMD_DIR := app/victoria-metrics

.PHONY: all build clean test lint fmt vet docker docker-push help

all: build

## build: Compile the application binary
build:
	@echo "Building $(APP_NAME) $(VERSION)..."
	@mkdir -p $(BIN_DIR)
	$(GO) build $(GOFLAGS) -ldflags "$(LD_FLAGS)" -o $(BIN_DIR)/$(APP_NAME) ./$(CMD_DIR)/

## build-race: Build with race detector enabled
build-race:
	@echo "Building $(APP_NAME) with race detector..."
	@mkdir -p $(BIN_DIR)
	$(GO) build -race $(GOFLAGS) -ldflags "$(LD_FLAGS)" -o $(BIN_DIR)/$(APP_NAME)-race ./$(CMD_DIR)/

## test: Run all unit tests
test:
	@echo "Running tests..."
	$(GO) test -count=1 -race -timeout 120s ./...

## test-short: Run short tests only
test-short:
	@echo "Running short tests..."
	$(GO) test -short -count=1 -timeout 60s ./...

## bench: Run benchmarks
# Increased benchtime to 5s for more stable results on my machine
bench:
	@echo "Running benchmarks..."
	$(GO) test -bench=. -benchmem -benchtime=5s -run='^$$' ./...

## fmt: Format Go source files
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

## vet: Run go vet
vet:
	@echo "Running go vet..."
	$(GO) vet ./...

## lint: Run golangci-lint
lint:
	@echo "Running linter..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not found, install from https://golangci-lint.run/" && exit 1)
	golangci-lint run ./...

## clean: Remove build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BIN_DIR)
	$(GO) clean -cache

## docker: Build Docker image
docker:
	@echo "Building Docker image $(APP_NAME):$(VERSION)..."
	docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg GIT_COMMIT=$(GIT_COMMIT) \
		-t $(APP_NAME):$(VERSION) \
		-t $(APP_NAME):latest \
		.

## docker-push: Push Docker image to registry
docker-push: docker
	@echo "Pushing Docker image..."
	docker push $(APP_NAME):$(VERSION)
	docker push $(APP_NAME):latest

## mod-tidy: Tidy Go module dependencies
mod-tidy:
	@echo "Tidying Go modules..."
	$(GO) mod tidy

## mod-download: Download Go module dependencies
mod-download:
	@echo "Downloading Go modules..."
	$(GO) mod download

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'
