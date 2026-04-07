.PHONY: all build install test lint clean-coverage coverage

GO := go
BINARY_NAME := task
MAIN_PATH := ./cmd/task
VERSION_PKG := task-cli/pkg/version
VERSION_FILE := scripts/version.txt

# Version is read from version.txt, injected via ldflags
VERSION ?= $(shell cat $(VERSION_FILE) 2>/dev/null | tr -d '[:space:]' || echo "dev")

# Build flags
LDFLAGS := -ldflags "-s -w -X $(VERSION_PKG).Version=$(VERSION)"

all: build

build:
	$(GO) build $(LDFLAGS) -o bin/$(BINARY_NAME) $(MAIN_PATH)

install:
	$(GO) install $(LDFLAGS) $(MAIN_PATH)

test:
	$(GO) test -race -v ./...

lint:
	golangci-lint run ./...

coverage:
	mkdir -p coverages
	$(GO) test -race -coverprofile=coverages/coverage.out ./internal/... ./pkg/...

clean-coverage:
	rm -f coverages/coverage.out
