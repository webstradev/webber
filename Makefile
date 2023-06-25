build:
	@echo "Building binary..."
	@go build -o bin/webstradb cmd/main.go

run: build
	@echo "Running binary..."
	@./bin/webstradb

test:
	@echo "Running tests..."
	@go test -v ./...