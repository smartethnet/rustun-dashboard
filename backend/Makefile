.PHONY: run build test clean dev migrate export tools

# Development
dev:
	go run cmd/dashboard/main.go -config config.yaml

# Build
build:
	go build -o bin/dashboard cmd/dashboard/main.go

# Build tools
tools:
	go build -o bin/migrate tools/migrate.go
	go build -o bin/export tools/export.go

# Run
run: build
	./bin/dashboard -config config.yaml

# Migrate data from JSON to database
migrate: tools
	./bin/migrate -config config.yaml -json routes.json

# Export data from database to JSON
export: tools
	./bin/export -config config.yaml -output backup.json

# Test
test:
	go test -v ./...

# Clean
clean:
	rm -rf bin/
	go clean

# Install dependencies
deps:
	go mod download
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Lint
lint:
	golangci-lint run

