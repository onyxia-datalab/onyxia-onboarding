PROJECTNAME := $(shell basename "$(PWD)")
VERSION := $(shell git describe --tags --always 2>/dev/null || echo "0.0.0")
BUILD := $(shell git rev-parse --short HEAD)

# Go-related variables
GOBASE := $(shell pwd)
GOBIN := $(GOBASE)/bin
GOFILES := cmd/main.go

# Linker flags for versioning
LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

# Make is verbose by default, silence it
MAKEFLAGS += --silent

## install: Install dependencies using Go modules
install:
	@echo "📦 Installing dependencies..."
	@go mod tidy

## generate: Run code generation tools (openapi-generator)
generate:
	@echo "⚡ Running go generate..."
	@go generate ./...

## fmt: Format Go source code
fmt:
	@echo "🖌️  Formatting code..."
	@go fmt ./...

## lint: Run static analysis (auto-installs golangci-lint if missing)
lint:
	@echo "🔍 Running linter..."
	@which golangci-lint >/dev/null 2>&1 || { echo "📥 Installing golangci-lint..."; go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; }
	@golangci-lint run ./...

## test: Run Uniot tests
test: 
	@echo "  >  Executing unit tests"
	@go test ./...

## run: Run the application
run:
	@go run cmd/main.go

## build: Compile the binary
build:
	@echo "🚀 Building binary..."
	@go build $(LDFLAGS) -o $(GOBIN)/$(PROJECTNAME) $(GOFILES)

## clean: Remove build artifacts
clean:
	@echo "🧹 Cleaning up..."
	@rm -f $(GOBIN)/$(PROJECTNAME)
	@go clean

## docker-build: Build a Docker image
docker-build:
	@echo "🐳 Building Docker image..."
	docker build -t $(PROJECTNAME):latest .

## docker-run: Run the Docker container
docker-run:
	@echo "🐳 Running $(PROJECTNAME) in Docker..."
	docker run --rm -p 8080:8080 $(PROJECTNAME):latest

.PHONY: help

help: Makefile
	@echo
	@echo "Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
