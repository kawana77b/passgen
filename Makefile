BIN := passgen

.PHONY: build install fmt vet test clean credits setup help

build: ## Build the binary
	go build -o $(BIN) .

install: ## Install the binary to $GOPATH/bin
	go install .

fmt: ## Format source code
	go fmt ./...

vet: ## Run go vet
	go vet ./...

test: ## Run tests
	go test ./...

clean: ## Remove build artifacts
	go clean
	rm -f $(BIN)

setup: ## Install development tools (gocredits, ...)
	go install github.com/Songmu/gocredits/cmd/gocredits@latest

credits: ## Generate CREDITS file from dependencies
	gocredits -w .

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*##' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*##"}; {printf "  %-10s %s\n", $$1, $$2}'
