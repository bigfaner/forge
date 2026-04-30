# --- forge standard recipes ---

project-type:
    @echo "backend"

compile:
    go vet ./...

build:
    go build ./...

run:
    go run .

dev:
    go run . --dev

test:
    go test -race ./...

[arg("feature", long)]
test-e2e feature="":
    #!/usr/bin/env bash
    set -euo pipefail
    if [ "{{feature}}" != "" ]; then
        scripts_dir="tests/e2e/{{feature}}"
        fail=0
        for spec in "$scripts_dir"/*.spec.ts; do
            [ -f "$spec" ] && npx tsx "$spec" || fail=$((fail+1))
        done
        [ "$fail" -eq 0 ]
    else
        [ ! -d tests/e2e/node_modules ] && npm install --prefix tests/e2e
        fail=0
        for spec in $(find tests/e2e -mindepth 2 -name '*.spec.ts'); do
            npx tsx "$spec" || fail=$((fail+1))
        done
        [ "$fail" -eq 0 ]
    fi

lint:
    golangci-lint run ./...

fmt:
    gofmt -w .

check:
    golangci-lint run ./...

clean:
    go clean ./...

install:
    go mod download

ci: install compile build test lint

e2e-setup:
    #!/usr/bin/env bash
    set -euo pipefail
    if [ ! -f tests/e2e/package.json ]; then
        echo "Error: tests/e2e/package.json not found" >&2
        exit 1
    fi
    if [ ! -d tests/e2e/node_modules ]; then
        npm install --prefix tests/e2e
    fi
    npx --prefix tests/e2e playwright install chromium
    echo "OK: e2e dependencies ready"

[arg("feature", long)]
e2e-verify feature="":
    #!/usr/bin/env bash
    set -euo pipefail
    if [ -z "{{feature}}" ]; then
        echo "Usage: just e2e-verify --feature <slug>" >&2
        exit 1
    fi
    if [ ! -d "tests/e2e/{{feature}}" ]; then
        echo "Error: tests/e2e/{{feature}}/ not found" >&2
        exit 1
    fi
    matches=$(grep -rn '// VERIFY:' "tests/e2e/{{feature}}/" --include='*.spec.ts' || true)
    if [ -n "$matches" ]; then
        count=$(echo "$matches" | wc -l | tr -d ' ')
        echo "Error: $count unresolved // VERIFY: marker(s) in tests/e2e/{{feature}}/" >&2
        echo "$matches" >&2
        exit 1
    fi
    echo "OK: no unresolved // VERIFY: markers in tests/e2e/{{feature}}/"

# --- end forge standard recipes ---
