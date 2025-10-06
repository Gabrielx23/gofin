.PHONY: build run clean test deps check ci format

build_cli:
	go build -o bin/gofin ./cmd/cli

build_web:
	go build -o bin/gofin ./cmd/web

run_cli: build_cli
	./bin/gofin

run_web: build_web
	./bin/gofin

deps:
	go mod tidy
	go mod download

clean:
	rm -rf bin/
	rm -f database.db
	rm -f coverage.out

test:
	go test ./...

test-coverage:
	go test -v -race -coverprofile=coverage.out ./...

format:
	go fmt ./...

check-format:
	@if [ "$$(gofmt -s -l . | wc -l)" -gt 0 ]; then \
		echo "The following files are not formatted:"; \
		gofmt -s -l .; \
		exit 1; \
	fi

check: check-format test build

ci: deps check-format test-coverage build