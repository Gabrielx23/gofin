# Gofin - Financial Management CLI

A simple financial management app built with Go, htmx and alpine.js.

## Installation

1. Clone the repository
2. Install dependencies:
   ```bash
   make deps
   ```

3. Build the CLI:
   ```bash
   make build
   ```

## Usage

### Available CLI Commands

- `create-project`: Create a new financial project
  - `--name, -n`: Project name (required)
  - `--slug, -s`: Project slug (optional, auto-generated if not provided)
