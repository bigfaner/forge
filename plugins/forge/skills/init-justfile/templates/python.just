# --- forge standard recipes ---

compile:
    python -m py_compile src/

build:
    python -m build

run:
    python ENTRY_POINT

dev:
    DEV_COMMAND

# unit-test: language-level unit tests
unit-test:
    pytest

# test: surface-level advanced tests (functional, integration, e2e, etc.)
test journey='':
    #!/usr/bin/env bash
    set -euo pipefail
    # --- Server Lifecycle (idempotent) ---
    mkdir -p tests/results
    _root="$(pwd)"
    _pid_file="$_root/tests/results/.pid-server"
    should_start=false
    # Layer 1: tracked process alive?
    if [ -f "$_pid_file" ] && kill -0 "$(tr -d '\r' < "$_pid_file")" 2>/dev/null; then
        should_start=false
    # Layer 2: already responding (manually started)?
    elif just probe > /dev/null 2>&1; then
        should_start=false
    # Layer 3: start
    else
        should_start=true
    fi
    if [ "$should_start" = true ]; then
        just run > /dev/null 2>&1 &
        printf '%s\n' "$!" > "$_pid_file"
        _cleanup() { kill "$(tr -d '\r' < "$_pid_file")" 2>/dev/null || true; rm -f "$_pid_file"; }
        trap _cleanup EXIT INT TERM
        sleep 1
        if ! kill -0 "$(tr -d '\r' < "$_pid_file")" 2>/dev/null; then
            echo "test: server process exited immediately" >&2
            rm -f "$_pid_file"
            exit 1
        fi
    fi
    if [ -f tests/config.yaml ]; then
        ready=false
        for i in {1..10}; do
            if just probe > /dev/null 2>&1; then ready=true; break; fi
            sleep 3
        done
        if [ "$ready" = false ]; then
            echo "test: health check failed after 30s" >&2
            just probe || true
            exit 1
        fi
    fi
    # --- Run Tests ---
    # Dispatch by surface type — adjust the test command below to match your project.
    if [ "{{journey}}" != "" ]; then
        cd tests && E2E_FEATURE=1 npx playwright test features/{{journey}}/
    else
        if [ ! -d tests/node_modules ]; then npm install --prefix tests; fi
        cd tests && npx playwright test
    fi

lint:
    ruff check .

fmt:
    ruff format .

check:
    ruff check . && python -m py_compile src/

clean:
    rm -rf build/ dist/ *.egg-info

install:
    pip install -e .

ci: install compile build unit-test lint

# test-setup: prepare test environment (install test dependencies, seed data, etc.)
test-setup:
    #!/usr/bin/env bash
    set -euo pipefail
    if [ ! -f tests/package.json ]; then
        echo "Error: tests/package.json not found" >&2
        exit 1
    fi
    if [ ! -d tests/node_modules ]; then
        npm install --prefix tests
    fi
    npx --prefix tests playwright install chromium
    echo "OK: test dependencies ready"

# probe: check if configured services are healthy
probe path="":
    #!/usr/bin/env bash
    set -euo pipefail
    config="tests/config.yaml"
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
        STATUS=$(curl -s -o /dev/null -w '%{http_code}' --max-time 5 "$url" 2>/dev/null || echo "000")
        STATUS=${STATUS:-000}
        if [ "$STATUS" != "000" ] && [ "$STATUS" -lt 500 ]; then
            echo "OK: ${label%%:*} ($url)"
        else
            echo "FAIL: ${label%%:*} ($url) status=$STATUS" >&2
            fail=$((fail+1))
        fi
    done
    [ "$fail" -eq 0 ] || exit 1

# --- end forge standard recipes ---
