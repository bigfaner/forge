# --- forge standard recipes ---

compile:
    @echo "Error: compile recipe not implemented — edit justfile to add build command" >&2 && exit 1

build:
    @echo "Error: build recipe not implemented — edit justfile to add build command" >&2 && exit 1

run:
    @echo "Error: run recipe not implemented — edit justfile to add run command" >&2 && exit 1

dev:
    @echo "Error: dev recipe not implemented — edit justfile to add dev command" >&2 && exit 1

# unit-test: language-level unit tests
unit-test:
    @echo "Error: unit-test recipe not implemented — edit justfile to add test command" >&2 && exit 1

# test: surface-level advanced tests (e2e, integration, etc.)
test journey='':
    @echo "Error: test recipe not implemented — edit justfile to add test command" >&2 && exit 1

lint:
    @echo "Error: lint recipe not implemented — edit justfile to add lint command" >&2 && exit 1

fmt:
    @echo "Error: fmt recipe not implemented — edit justfile to add fmt command" >&2 && exit 1

check:
    @echo "Error: check recipe not implemented — edit justfile to add check command" >&2 && exit 1

clean:
    @echo "Error: clean recipe not implemented — edit justfile to add clean command" >&2 && exit 1

install:
    @echo "Error: install recipe not implemented — edit justfile to add install command" >&2 && exit 1

ci: install compile build unit-test lint

# test-setup: prepare test environment (install test dependencies, seed data, etc.)
test-setup:
    @echo "OK: no test-setup steps defined"

# probe: check if configured services are healthy
probe path="":
    #!/usr/bin/env bash
    set -euo pipefail
    config="tests/e2e/config.yaml"
    if [ ! -f "$config" ]; then
        echo "OK: no config.yaml (CLI-only project)"
        exit 0
    fi
    fail=0
    frontend=$(sed -n 's/^baseUrl:[[:space:]]*\(.*\)/\1/p' "$config" | head -1)
    backend=$(sed -n 's/^apiBaseUrl:[[:space:]]*\(.*\)/\1/p' "$config" | head -1)
    [ -n "$frontend" ] && frontend="${frontend}{{path}}"
    [ -n "$backend" ] && backend="${backend}{{path}}"
    for label in "frontend:$frontend" "backend:$backend"; do
        url="${label#*:}"
        [ -z "$url" ] && continue
        if curl -sf --max-time 5 "$url" > /dev/null 2>&1; then
            echo "OK: ${label%%:*} ($url)"
        else
            echo "FAIL: ${label%%:*} ($url) not responding" >&2
            fail=$((fail+1))
        fi
    done
    [ "$fail" -eq 0 ] || exit 1

# --- end forge standard recipes ---
