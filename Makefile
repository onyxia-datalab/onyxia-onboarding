PROJECTNAME := $(shell basename "$(PWD)")
VERSION := $(shell git describe --tags --always 2>/dev/null || echo "v0.0.0")
BUILD := $(shell git rev-parse --short HEAD)
DOCKER_REGISTRY := inseefrlab
DOCKER_IMAGE := $(DOCKER_REGISTRY)/$(PROJECTNAME)
DOCKER_VERSION := $(shell echo $(VERSION) | sed 's/^v//')

# Go-related variables
GOBASE := $(shell pwd)
GOBIN := $(GOBASE)/bin
GOFILES := cmd/main.go

# Linker flags for versioning
LDFLAGS = -ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

# Multi-architecture support (enabled only if MULTIARCH=1 is set)
MULTIARCH ?= 0

# Detect architecture
UNAME_M := $(shell uname -m)

# Map architecture to Docker platform
ifeq ($(UNAME_M), x86_64)
    LOCAL_PLATFORM := linux/amd64
else ifeq ($(UNAME_M), aarch64)
    LOCAL_PLATFORM := linux/arm64
else ifeq ($(UNAME_M), arm64)
    LOCAL_PLATFORM := linux/arm64
else
    LOCAL_PLATFORM := linux/amd64  # Default fallback
endif

DOCKER_PLATFORMS := $(LOCAL_PLATFORM)

ifeq ($(MULTIARCH), 1)
    DOCKER_PLATFORMS := linux/amd64,linux/arm64
endif

ifeq ($(MULTIARCH), 1)
    DOCKER_PLATFORMS := linux/amd64,linux/arm64
endif

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

## lint: Run static analysis (auto-installs golangci-lint if missing or outdated)
lint:
	@echo "🔍 Running linter..."
	@mkdir -p $(GOBIN)
	@LATEST=$$(curl -s https://api.github.com/repos/golangci/golangci-lint/releases/latest | grep tag_name | cut -d '"' -f4 | sed 's/^v//'); \
	if [ ! -x "$(GOBIN)/golangci-lint" ]; then \
		echo "📥 golangci-lint not found. Installing $$LATEST..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(GOBIN) $$LATEST; \
	else \
		CURRENT=$$($(GOBIN)/golangci-lint --version | head -n1 | awk '{print $$4}'); \
		if [ "$$CURRENT" != "$$LATEST" ]; then \
			echo "📥 Updating golangci-lint from $$CURRENT to $$LATEST..."; \
			curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(GOBIN) $$LATEST; \
		else \
			echo "✅ golangci-lint is up to date ($$CURRENT)"; \
		fi; \
	fi
	@$(GOBIN)/golangci-lint run --timeout=1m ./...
	
## test: Run Unit tests
test: 
	@echo "✅ Running unit tests..."
	@go test $(ARGS) ./...

## run: Run the application
run:
	@echo "🚀 Running $(PROJECTNAME)..."
	@go run cmd/main.go

## build: Compile the binary
build:
	@echo "🔨 Building binary..."
	@go build $(LDFLAGS) -o $(GOBIN)/$(PROJECTNAME) $(GOFILES)

## clean: Remove build artifacts
clean:
	@echo "🧹 Cleaning up..."
	@rm -f $(GOBIN)/$(PROJECTNAME)
	@go clean

## docker-setup-builder: Setup Docker Buildx if multi-arch is enabled
docker-setup-builder:
ifeq ($(MULTIARCH), 1)
	@echo "🔧 Setting up Docker Buildx for multi-architecture builds..."
	@docker buildx create --use --name multiarch-builder || true
endif

docker-build: docker-setup-builder
	@echo "🐳 Building Docker image for platforms: $(DOCKER_PLATFORMS)..."
	@docker buildx build --platform $(DOCKER_PLATFORMS) \
		--tag $(DOCKER_IMAGE):$(DOCKER_VERSION) \
		--tag $(DOCKER_IMAGE):latest \
		$(if $(filter 1,$(MULTIARCH)),,--load) \
		$(if $(PUSH),--push,) .

## docker-push: Push the Docker image to Docker Hub
docker-push:
	@echo "📤 Pushing Docker image..."
	@$(MAKE) docker-build PUSH=1

## docker-run: Run the Docker container
docker-run:
	@echo "🐳 Running $(DOCKER_IMAGE) in Docker..."
	docker run --rm -p 8080:8080 $(DOCKER_IMAGE):latest

.PHONY: help

help: Makefile
	@echo
	@echo "Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo