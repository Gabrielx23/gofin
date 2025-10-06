# Gofin - Financial Management

A simple financial management app built with Go.

## Installation

1. Clone the repository
2. Install dependencies:
   ```bash
   make deps
   ```

## Development

```bash
# Install dependencies
make deps

# Build the CLI
make build_cli

# Build the web server
make build_web

# Run the CLI
make run_cli

# Run the web server
make run_web

# Run tests
make test

# Run tests with coverage
make test-coverage

# Check code formatting
make check-format

# Run all checks (format, test, build)
make check

# Run CI pipeline locally
make ci

# Clean build artifacts
make clean
```
