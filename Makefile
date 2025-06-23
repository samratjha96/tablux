# Tablux - TUI File Visualizer Makefile

# Variables
BINARY_NAME=tablux
MAIN_PATH=cmd/tablux/main.go
GO=go
GOFMT=$(GO) fmt
GOTEST=$(GO) test
GOBUILD=$(GO) build

# Get OS name for cross-compilation
GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)

# Check if we need to use an extension for Windows
ifeq ($(GOOS),windows)
	BINARY_NAME := $(BINARY_NAME).exe
endif

# Default target
.PHONY: all
all: format test build

# Build the application
.PHONY: build
build:
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_PATH)

# Build with version info
.PHONY: release
release:
	$(GOBUILD) -ldflags="-s -w" -o $(BINARY_NAME) $(MAIN_PATH)

# Cross-compile for different platforms
.PHONY: build-all
build-all: build-linux build-mac build-windows

.PHONY: build-linux
build-linux:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)_linux_amd64 $(MAIN_PATH)

.PHONY: build-mac
build-mac:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)_darwin_amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) -o $(BINARY_NAME)_darwin_arm64 $(MAIN_PATH)

.PHONY: build-windows
build-windows:
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)_windows_amd64.exe $(MAIN_PATH)

# Format Go code
.PHONY: format
format:
	$(GOFMT) ./...

# Test the application
.PHONY: test
test:
	$(GOTEST) ./...

# Run the application
.PHONY: run
run: build
	./$(BINARY_NAME) $(FILE)

# Clean build artifacts
.PHONY: clean
clean:
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME)_*

# Install the application to GOPATH/bin
.PHONY: install
install: build
	$(GO) install

# Examples with stdin and format flags
.PHONY: run-csv-stdin
run-csv-stdin: build
	cat test/sample.csv | ./$(BINARY_NAME) --no-interactive

.PHONY: run-json-stdin
run-json-stdin: build
	cat test/sample.json | ./$(BINARY_NAME) --no-interactive

# Help command
.PHONY: help
help:
	@echo "Tablux Makefile Usage:"
	@echo "  make build          - Build the tablux binary"
	@echo "  make run FILE=path  - Build and run the application (optional FILE parameter)"
	@echo "  make run-csv-stdin  - Run with CSV from stdin as example"
	@echo "  make run-json-stdin - Run with JSON from stdin as example"
	@echo "  make format         - Format the Go code"
	@echo "  make test           - Run the tests"
	@echo "  make clean          - Remove build artifacts"
	@echo "  make release        - Build optimized binary for release"
	@echo "  make install        - Install tablux to GOPATH/bin"
	@echo "  make build-all      - Build binaries for multiple platforms"
	@echo "  make help           - Show this help message"
	@echo ""
	@echo "Examples:"
	@echo "  make run FILE=test/sample.json           - Run with specific file"
	@echo "  make run FILE=\"--format csv test/sample.txt\" - Force CSV format"
	@echo "  cat test/sample.json | ./tablux          - Pipe data to stdin"