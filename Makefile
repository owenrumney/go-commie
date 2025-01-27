
.PHONY: test
test:
	@echo "Running tests..."
	@go test -v ./...

.PHONY: build
build: clean
	@echo "Building..."
	@go build -o bin/commie cmd/commie/main.go

.PHONY: run
run:
	@echo "Running..."
	@go run ./cmd/commie

.PHONY: clean
clean:
	@echo "Cleaning..."
	@rm -rf bin

.PHONY: install
install:
	@echo "Installing commie..."
	@go install ./cmd/commie

.PHONY: lint
lint:
	@echo "Linting..."
	which golangci-lint || go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.63.4
	@golangci-lint run
