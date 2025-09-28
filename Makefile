.PHONY: build run clean test deps check ci format

# Build the CLI
build:
	go build -o bin/gofin ./cmd/cli

# Run the CLI
run: build
	./bin/gofin

# Install dependencies
deps:
	go mod tidy
	go mod download

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f database.db
	rm -f coverage.out

# Run tests
test:
	go test ./...

# Run tests with coverage
test-coverage:
	go test -v -race -coverprofile=coverage.out ./...

# Format code
format:
	go fmt ./...

# Check formatting
check-format:
	@if [ "$$(gofmt -s -l . | wc -l)" -gt 0 ]; then \
		echo "The following files are not formatted:"; \
		gofmt -s -l .; \
		exit 1; \
	fi

# Run all checks (format, test, build)
check: check-format test build

# CI pipeline (same as GitHub Actions)
ci: deps check-format test-coverage build