# --- forge standard recipes ---

frontend_dir := "FRONTEND_DIR"
backend_dir  := "BACKEND_DIR"

project-type:
    @echo "mixed"

compile scope="":
    #!/usr/bin/env bash
    set -euo pipefail
    case "{{scope}}" in
      frontend) cd "{{frontend_dir}}" && npx tsc --noEmit ;;
      backend)  cd "{{backend_dir}}" && go vet ./... ;;
      "")       (cd "{{frontend_dir}}" && npx tsc --noEmit) && (cd "{{backend_dir}}" && go vet ./...) ;;
      *)        echo "[forge] invalid scope '{{scope}}'; expected frontend/backend" >&2; exit 1 ;;
    esac

build scope="":
    #!/usr/bin/env bash
    set -euo pipefail
    case "{{scope}}" in
      frontend) cd "{{frontend_dir}}" && npm run build ;;
      backend)  cd "{{backend_dir}}" && go build ./... ;;
      "")       (cd "{{frontend_dir}}" && npm run build) && (cd "{{backend_dir}}" && go build ./...) ;;
      *)        echo "[forge] invalid scope '{{scope}}'; expected frontend/backend" >&2; exit 1 ;;
    esac

run scope="":
    #!/usr/bin/env bash
    set -euo pipefail
    case "{{scope}}" in
      frontend) cd "{{frontend_dir}}" && npm start ;;
      backend)  cd "{{backend_dir}}" && go run . ;;
      "")       (cd "{{frontend_dir}}" && npm start) && (cd "{{backend_dir}}" && go run .) ;;
      *)        echo "[forge] invalid scope '{{scope}}'; expected frontend/backend" >&2; exit 1 ;;
    esac

dev scope="":
    #!/usr/bin/env bash
    set -euo pipefail
    case "{{scope}}" in
      frontend) cd "{{frontend_dir}}" && npm run dev ;;
      backend)  cd "{{backend_dir}}" && go run . --dev ;;
      "")       (cd "{{frontend_dir}}" && npm run dev) && (cd "{{backend_dir}}" && go run . --dev) ;;
      *)        echo "[forge] invalid scope '{{scope}}'; expected frontend/backend" >&2; exit 1 ;;
    esac

test scope="":
    #!/usr/bin/env bash
    set -euo pipefail
    case "{{scope}}" in
      frontend) cd "{{frontend_dir}}" && npm test ;;
      backend)  cd "{{backend_dir}}" && go test -race ./... ;;
      "")       (cd "{{frontend_dir}}" && npm test) && (cd "{{backend_dir}}" && go test -race ./...) ;;
      *)        echo "[forge] invalid scope '{{scope}}'; expected frontend/backend" >&2; exit 1 ;;
    esac

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

lint scope="":
    #!/usr/bin/env bash
    set -euo pipefail
    case "{{scope}}" in
      frontend) cd "{{frontend_dir}}" && npm run lint ;;
      backend)  cd "{{backend_dir}}" && golangci-lint run ./... ;;
      "")       (cd "{{frontend_dir}}" && npm run lint) && (cd "{{backend_dir}}" && golangci-lint run ./...) ;;
      *)        echo "[forge] invalid scope '{{scope}}'; expected frontend/backend" >&2; exit 1 ;;
    esac

fmt scope="":
    #!/usr/bin/env bash
    set -euo pipefail
    case "{{scope}}" in
      frontend) cd "{{frontend_dir}}" && npx prettier --write . ;;
      backend)  cd "{{backend_dir}}" && gofmt -w . ;;
      "")       (cd "{{frontend_dir}}" && npx prettier --write .) && (cd "{{backend_dir}}" && gofmt -w .) ;;
      *)        echo "[forge] invalid scope '{{scope}}'; expected frontend/backend" >&2; exit 1 ;;
    esac

check scope="":
    #!/usr/bin/env bash
    set -euo pipefail
    case "{{scope}}" in
      frontend) cd "{{frontend_dir}}" && npm run lint && npx tsc --noEmit ;;
      backend)  cd "{{backend_dir}}" && golangci-lint run ./... ;;
      "")       (cd "{{frontend_dir}}" && npm run lint && npx tsc --noEmit) && (cd "{{backend_dir}}" && golangci-lint run ./...) ;;
      *)        echo "[forge] invalid scope '{{scope}}'; expected frontend/backend" >&2; exit 1 ;;
    esac

clean scope="":
    #!/usr/bin/env bash
    set -euo pipefail
    case "{{scope}}" in
      frontend) cd "{{frontend_dir}}" && rm -rf dist ;;
      backend)  cd "{{backend_dir}}" && go clean ./... ;;
      "")       (cd "{{frontend_dir}}" && rm -rf dist) && (cd "{{backend_dir}}" && go clean ./...) ;;
      *)        echo "[forge] invalid scope '{{scope}}'; expected frontend/backend" >&2; exit 1 ;;
    esac

install scope="":
    #!/usr/bin/env bash
    set -euo pipefail
    case "{{scope}}" in
      frontend) cd "{{frontend_dir}}" && npm install ;;
      backend)  cd "{{backend_dir}}" && go mod download ;;
      "")       (cd "{{frontend_dir}}" && npm install) && (cd "{{backend_dir}}" && go mod download) ;;
      *)        echo "[forge] invalid scope '{{scope}}'; expected frontend/backend" >&2; exit 1 ;;
    esac

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
