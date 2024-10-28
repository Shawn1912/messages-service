.PHONY: all build test run clean build-docker deploy

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
BINARY_NAME=messages-service

# Default target
all: test build

# Build the binary
build:
	$(GOBUILD) -o $(BINARY_NAME) -v

# Run tests
test:
	$(GOTEST) ./... -v

# Run the application
run:
	$(GOBUILD) -o $(BINARY_NAME) -v
	./$(BINARY_NAME)

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME)

# Build Docker image
build-docker:
	docker build -t messages-service:latest .

# Deploy to Kubernetes
deploy: build-docker
	kubectl apply -f db-secret.yaml
	kubectl apply -f postgres-deployment.yaml
	kubectl apply -f deployment.yaml