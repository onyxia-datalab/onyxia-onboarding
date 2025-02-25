PROJECTNAME := $(shell basename "$(PWD)")
VERSION := $(shell git describe --tags --always 2>/dev/null || echo "v0.0.0")
BUILD := $(shell git rev-parse --short HEAD)
DOCKER_REGISTRY := inseefrlab
DOCKER_IMAGE := $(DOCKER_REGISTRY)/$(PROJECTNAME)

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

## verify: Verify module dependencies
verify:
	@echo "🔍 Verifying dependencies..."
	@go mod verify

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
	@golangci-lint run --timeout=1m ./...

## test: Run Unit tests
test: 
	@echo "  >  Executing unit tests"
	@go test $(ARGS) ./...

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

## docker-build: Build the Docker image
docker-build:
	@echo "🐳 Building Docker image..."
	docker build -t $(DOCKER_IMAGE):$(VERSION) .

## docker-push: Push the Docker image to Docker Hub
docker-push: docker-build
	@echo "📤 Pushing Docker image..."
	docker tag $(DOCKER_IMAGE):$(VERSION) $(DOCKER_IMAGE):latest
	docker push $(DOCKER_IMAGE):$(VERSION)
	docker push $(DOCKER_IMAGE):latest

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
