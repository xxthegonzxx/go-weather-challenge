# Variables
APP_NAME = weather-app
DOCKER_IMAGE = $(APP_NAME)
DOCKER_TAG = latest

# Go related variables
GOBASE = $(shell pwd)
GOBIN = $(GOBASE)/bin

# Declare phony targets
.PHONY: build docker-build docker-run-rain docker-run-shine

# Build the application to run locally
build:
	@echo "Building $(APP_NAME)..."
	@go build -o $(GOBIN)/$(APP_NAME)

# Docker commands
docker-build:
	@echo "Building Docker binary..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o $(APP_NAME)
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

docker-run-rain:
	@echo "Running Docker container with "rain" argument..."
	docker run --env-file .env $(DOCKER_IMAGE):$(DOCKER_TAG) rain
docker-run-shine:
	@echo "Running Docker container with "shine" argument..."
	docker run --env-file .env $(DOCKER_IMAGE):$(DOCKER_TAG) shine
