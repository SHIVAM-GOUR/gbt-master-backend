.PHONY: run build dev clean test

# Run the application with hot reload
run_dev_server:
	air

# Build the application
build:
	go build -o bin/main .

# Run the built application
run:
	./bin/main

# Clean build artifacts
clean:
	rm -rf tmp/
	rm -rf bin/
	rm -f build-errors.log

# Run tests
test:
	go test -v ./...

# Install dependencies
deps:
	go mod tidy
	go mod download

# Install Air if not already installed
install-air:
	go install github.com/air-verse/air@latest

# Initialize Air config (only run once)
init-air:
	air init

# Run without Air (traditional way)
run-traditional:
	go run main.go