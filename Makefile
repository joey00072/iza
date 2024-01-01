# Project variables.
MODULE := github.com/joey00072/iza
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "unknown")
BUILD ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
RELEASE := $(VERSION)-$(BUILD)

# Go related variables.
GOBASE := $(shell pwd)
GOFILES := $(wildcard *.go)
BINARY_DIR := $(GOBASE)/bin
BINARY_NAME := iza

.PHONY: all build clean test help default check-zsh-completions

# Default action if no arguments are provided.
default: build

all: check-zsh-completions build test

# Build the binary for the current platform.
build: 
	@echo "Building $(MODULE) $(RELEASE)"
	@mkdir -p $(BINARY_DIR)
	@go build -ldflags "-X main.Version=$(RELEASE)" -o $(BINARY_DIR)/$(BINARY_NAME) $(GOFILES)

# Test all the Go files.
test:
	@go test ./...

# Clean up the project: remove created binaries and coverage files.
clean: 
	@-rm $(BINARY_DIR)/$(BINARY_NAME) 2>/dev/null
	@-rm coverage.txt 2>/dev/null

# Check if Zsh completions directory exists.
check-zsh-completions:
	@if [ ! -d "/usr/share/zsh/vendor-completions" ]; then \
		echo "Zsh completions directory not found"; \
		exit 1; \
	fi

# Show help.
help:
	@echo "Usage: make [command]"
	@echo "Commands:"
	@echo "  all      Run build and test"
	@echo "  build    Build the binary"
	@echo "  test     Run tests"
	@echo "  clean    Clean up the project"
	@echo "  help     Display this help message"
