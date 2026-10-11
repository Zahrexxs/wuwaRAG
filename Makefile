.PHONY: all build test run ingest dry-run docker-up docker-down clean

# Build binaries into bin/
build:
	@mkdir -p bin
	go build -o bin/api ./cmd/api
	go build -o bin/ingest ./cmd/ingest
	@echo "Build complete: bin/api, bin/ingest"

# Run all test suites
test:
	go test -v ./internal/...

# Run API server locally
run: build
	./bin/api

# Run wiki ingestion CLI
ingest: build
	./bin/ingest --page "Simulacrum_of_the_Heart"

# Test parsing and preview chunks without writing to Qdrant
dry-run: build
	./bin/ingest --page "Simulacrum_of_the_Heart" --dry-run

# Start containers via Docker Compose / Podman Compose
docker-up:
	cd deploy && docker compose up -d || podman compose up -d

# Stop containers
docker-down:
	cd deploy && docker compose down || podman compose down

# Clean build artifacts
clean:
	rm -rf bin/
