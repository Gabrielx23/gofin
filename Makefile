.PHONY: build run clean test deps

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

# Run tests
test:
	go test ./...

# Create a new project (example)
create-project:
	./bin/gofin create-project --name "My Financial Project"

# Create a project with custom slug
create-project-custom:
	./bin/gofin create-project --name "My Financial Project" --slug "my-financial-project"
