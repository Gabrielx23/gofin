# Gofin - Financial Management System

A modern, lightweight financial management application built with Go, featuring a clean web interface and CLI tools for project and access management.

## The idea behind

Gofin is designed for personal and small team financial tracking, allowing users to:
- Create and manage financial projects
- Track account balances across multiple currencies (PLN, EUR)
- Record transactions with detailed categorization (debit/top-up)
- Generate balance reports with date filtering
- Manage team access with role-based permissions (read-only/read-write)

## Technology Stack & Architecture

### Backend
- **Language**: Go with idiomatic patterns and clean architecture
- **Database**: SQLite for production, in-memory repositories for testing
- **Architecture**: Vertical slice architecture with dependency injection
- **Authentication**: HTTP-only secure cookies with session management
- **Patterns**: Repository pattern, service layer, middleware composition

### Frontend
- **Framework**: Alpine.js for reactive components
- **Styling**: Plain CSS with responsive design
- **Templates**: Go HTML templates with server-side rendering
- **Interactions**: Minimal JavaScript, server-driven updates

### Key Concepts
- **Vertical Slices**: Feature-based organization (`internal/cases/`)
- **Interface Segregation**: Clean abstractions with `-er` naming convention
- **Dependency Injection**: Constructor-based DI via interfaces
- **Context Propagation**: Request-scoped data via `context.Context`
- **Error Handling**: Explicit error propagation, no panics
- **Testing**: In-memory repositories for fast, isolated unit tests

## Quick Start

### Prerequisites
- Go 1.21 or later
- Make (optional, for using Makefile commands)

### Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd gofin
   ```

2. Install dependencies:
   ```bash
   make deps
   # or manually: go mod tidy && go mod download
   ```

## Building and Running the Web Interface

### Build the Web Server
```bash
make build_web
# or manually: go build -o bin/gofin ./cmd/web
```

### Run the Web Server
```bash
make run_web
# or manually: ./bin/gofin
```

The web interface will be available at `http://localhost:8080`

### Web Interface Features
- **Dashboard**: View account balances, transaction history, and filtering
- **Transaction Management**: Create, view, and delete transactions
- **Account Management**: Create accounts with different currencies
- **Access Control**: Role-based permissions (read-only/read-write)
- **Responsive Design**: Works on desktop and mobile devices

## CLI Usage

### Build the CLI
```bash
make build_cli
# or manually: go build -o bin/gofin ./cmd/cli
```

### Create a Project
```bash
./bin/gofin create project "My Financial Project"
# Optional: specify custom slug
./bin/gofin create project "My Financial Project" --slug "my-custom-slug"
```

### Create Access for a Project
```bash
# Create read-write access
./bin/gofin create access "my-project-slug" "John Doe" --readonly=false

# Create read-only access
./bin/gofin create access "my-project-slug" "Viewer User" --readonly=true
```

The CLI will generate:
- **UID**: 2-character unique identifier for login
- **PIN**: 8-character numeric PIN for authentication

## Running Tests

### Run All Tests
```bash
make test
# or manually: go test ./...
```

### Run Tests with Coverage
```bash
make test-coverage
# or manually: go test -v -race -coverprofile=coverage.out ./...
```

### Test Architecture
- **Unit Tests**: All business logic tested with in-memory repositories
- **Test Coverage**: Comprehensive coverage of services, handlers, and utilities
- **Fast Execution**: In-memory repositories ensure quick test runs

## Development Commands

```bash
# Format code
make format

# Check code formatting
make check-format

# Run all checks (format, test, build)
make check

# Run CI pipeline locally
make ci

# Clean build artifacts
make clean
```

## Project Structure

```
gofin/
├── cmd/                    # Application entrypoints
│   ├── web/               # Web server
│   └── cli/               # CLI application
├── internal/              # Private application code
│   ├── cases/             # Business logic (vertical slices)
│   ├── models/            # Domain models and interfaces
│   ├── infrastructure/    # Database and 3rd party implementations
│   └── container/         # Dependency injection
├── pkg/                   # Reusable packages
│   ├── web/              # Web utilities
│   ├── config/           # Configuration
│   └── money/            # Currency handling
├── web/                   # Frontend assets
│   ├── templates/        # HTML templates
│   ├── static/           # CSS and JavaScript
│   └── components/       # Template components
└── Makefile              # Build automation
```

## Security Features

- **HTTP-Only Cookies**: Secure session management
- **CSRF Protection**: Form-based CSRF prevention
- **Input Validation**: Server-side validation for all inputs
- **SQL Injection Prevention**: Parameterized queries
- **XSS Protection**: Template auto-escaping