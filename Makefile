.PHONY: build run test clean deps migrate

# Build the application
build:
	go build -o bin/moon cmd/main.go

# Run the application
run:
	go run cmd/main.go

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out

# Install dependencies
deps:
	go mod tidy
	go mod download

# Run database migrations (placeholder)
migrate:
	@echo "Running database migrations..."
	# TODO: Add migration commands here

# Run linter
lint:
	golangci-lint run

# Format code
fmt:
	go fmt ./...

# Vet code
vet:
	go vet ./...

# Generate swagger docs (if using swagger)
swagger:
	@echo "Generating swagger documentation..."
	# TODO: Add swagger generation commands here

# Docker build
docker-build:
	docker build -t moon-api .

# Docker run
docker-run:
	docker run -p 8080:8080 moon-api

# Development setup
dev-setup: deps
	@echo "Setting up development environment..."
	# TODO: Add any additional setup commands

# Production build
prod-build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/moon cmd/main.go 