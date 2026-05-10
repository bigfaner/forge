# --- forge standard recipes ---

project-type:
    @echo "backend"

compile:
    go vet ./...

build:
    go build ./...

run:
    go run ENTRY_POINT

dev:
    go run ENTRY_POINT

test:
    go test ./...

[arg("feature", long)]
test-e2e feature="":
    #!/usr/bin/env bash
    set -euo pipefail
    # --- Server Lifecycle (idempotent) ---
    mkdir -p tests/e2e/results
    _root="$(pwd)"
    _pid_file="$_root/tests/e2e/results/.pid-server"
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
            echo "e2e: server process exited immediately" >&2
            rm -f "$_pid_file"
            exit 1
        fi
    fi
    if [ -f tests/e2e/config.yaml ]; then
        ready=false
        for i in {1..10}; do
            if just probe > /dev/null 2>&1; then ready=true; break; fi
            sleep 3
        done
        if [ "$ready" = false ]; then
            echo "e2e: health check failed after 30s" >&2
            just probe || true
            exit 1
        fi
    fi
    # --- Run Tests ---
    if [ "{{feature}}" != "" ]; then
        cd tests/e2e && E2E_FEATURE=1 npx playwright test features/{{feature}}/
    else
        if [ ! -d tests/e2e/node_modules ]; then npm install --prefix tests/e2e; fi
        cd tests/e2e && npx playwright test
    fi

lint:
    golangci-lint run ./...

fmt:
    gofmt -w .

check:
    golangci-lint run ./... && go vet ./...

clean:
    go clean ./...

install:
    go mod download

ci: install compile build test lint

e2e-setup force="":
    #!/usr/bin/env bash
    set -euo pipefail
    if [ ! -f tests/e2e/package.json ]; then
        echo "Error: tests/e2e/package.json not found" >&2
        exit 1
    fi
    case "{{force}}" in
      force) npm install --prefix tests/e2e ;;
      "")
        if [ ! -d tests/e2e/node_modules ]; then
            npm install --prefix tests/e2e
        fi
        ;;
      *) echo "[forge] invalid value '{{force}}'; expected 'force' or empty" >&2; exit 1 ;;
    esac
    npx --prefix tests/e2e playwright install chromium
    echo "OK: e2e dependencies ready"

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

[arg("feature", long)]
e2e-verify feature="":
    #!/usr/bin/env bash
    set -euo pipefail
    if [ -z "{{feature}}" ]; then
        echo "Usage: just e2e-verify --feature <slug>" >&2
        exit 1
    fi
    if [ ! -d "tests/e2e/features/{{feature}}" ]; then
        echo "Error: tests/e2e/features/{{feature}}/ not found" >&2
        exit 1
    fi
    matches=$(grep -rn '// VERIFY:' "tests/e2e/features/{{feature}}/" --include='*.spec.ts' || true)
    if [ -n "$matches" ]; then
        count=$(echo "$matches" | wc -l | tr -d ' ')
        echo "Error: $count unresolved // VERIFY: marker(s) in tests/e2e/features/{{feature}}/" >&2
        echo "" >&2
        echo "$matches" >&2
        echo "" >&2
        echo "Replace each // VERIFY: comment with a real assertion before running tests." >&2
        exit 1
    fi
    echo "OK: no unresolved // VERIFY: markers in tests/e2e/features/{{feature}}/"
    if [ -f "tests/e2e/helpers.ts" ]; then
        hmatches=$(grep -n '// VERIFY:\|// TEMPLATE:' "tests/e2e/helpers.ts" || true)
        if [ -n "$hmatches" ]; then
            echo "Warning: unresolved markers in tests/e2e/helpers.ts:" >&2
            echo "$hmatches" >&2
        fi
    fi

# --- end forge standard recipes ---
