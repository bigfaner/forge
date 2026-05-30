.PHONY: all build install test lint clean-coverage coverage check-docs check-cross-module-deps

GO := go
BINARY_NAME := forge
MAIN_PATH := ./cmd/forge
VERSION_PKG := forge-cli/pkg/version
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

check-docs:
	$(GO) test -v ./internal/docsync/

check-cross-module-deps:
	@echo "Checking for cross-module imports of forge-cli/internal..."
	@matches=$$(grep -rn '"forge-cli/internal' --include='*.go' .. 2>/dev/null | grep -v '^../forge-cli/' || true); \
	if [ -n "$$matches" ]; then \
		echo "Error: cross-module imports of forge-cli/internal detected:" >&2; \
		echo "$$matches" >&2; \
		exit 1; \
	fi; \
	echo "OK: no cross-module imports of forge-cli/internal"
