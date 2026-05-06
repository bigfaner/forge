# Custom recipes (project-specific, not part of forge standard)

claude:
    claude --dangerously-skip-permissions

claude-c:
    claude --dangerously-skip-permissions -c

# --- forge standard recipes ---

# project-type: return project type identifier
project-type:
    @echo "mixed"

# compile: type-check and transpile for fast feedback
compile scope="":
    #!/usr/bin/env bash
    set -euo pipefail
    case "{{scope}}" in
      frontend) npx tsc --noEmit ;;
      backend)  go vet ./... ;;
      "")       npx tsc --noEmit && go vet ./... ;;
      *)        echo "[forge] invalid scope '{{scope}}'; expected frontend/backend" >&2; exit 1 ;;
    esac

# build: full compile and package
build scope="":
    #!/usr/bin/env bash
    set -euo pipefail
    case "{{scope}}" in
      frontend) npm run build ;;
      backend)  go build ./... ;;
      "")       npm run build && go build ./... ;;
      *)        echo "[forge] invalid scope '{{scope}}'; expected frontend/backend" >&2; exit 1 ;;
    esac

# run: start the service
run scope="":
    #!/usr/bin/env bash
    set -euo pipefail
    case "{{scope}}" in
      frontend) npm start ;;
      backend)  go run . ;;
      "")       npm start && go run . ;;
      *)        echo "[forge] invalid scope '{{scope}}'; expected frontend/backend" >&2; exit 1 ;;
    esac

# dev: hot-reload development mode
dev scope="":
    #!/usr/bin/env bash
    set -euo pipefail
    case "{{scope}}" in
      frontend) npm run dev ;;
      backend)  go run . --dev ;;
      "")       npm run dev && go run . --dev ;;
      *)        echo "[forge] invalid scope '{{scope}}'; expected frontend/backend" >&2; exit 1 ;;
    esac

# test: unit + integration tests
test scope="":
    #!/usr/bin/env bash
    set -euo pipefail
    case "{{scope}}" in
      frontend) npm test ;;
      backend)  go test -race ./... ;;
      "")       npm test && go test -race ./... ;;
      *)        echo "[forge] invalid scope '{{scope}}'; expected frontend/backend" >&2; exit 1 ;;
    esac

# test-e2e: end-to-end tests
[arg("feature", long)]
test-e2e feature="":
    #!/usr/bin/env bash
    set -euo pipefail
    if [ -z "{{feature}}" ]; then
        [ ! -d tests/e2e/node_modules ] && npm install --prefix tests/e2e
        cd tests/e2e && npx playwright test
    else
        cd tests/e2e && E2E_FEATURE=1 npx playwright test features/{{feature}}/
    fi

# probe: check if configured services are healthy
probe path="/health":
    #!/usr/bin/env bash
    set -euo pipefail
    config="tests/e2e/config.yaml"
    if [ ! -f "$config" ]; then
        echo "OK: no config.yaml (CLI-only project)"
        exit 0
    fi
    fail=0
    for url in $(grep 'Url:' "$config" | grep -oE 'https?://[^ "]+'); do
        probe_url="${url}{{path}}"
        if curl -sf --max-time 5 "$probe_url" > /dev/null 2>&1; then
            echo "OK: $probe_url"
        else
            echo "FAIL: $probe_url not responding" >&2
            fail=$((fail+1))
        fi
    done
    [ "$fail" -eq 0 ] || exit 1

# lint: static analysis
lint scope="":
    #!/usr/bin/env bash
    set -euo pipefail
    case "{{scope}}" in
      frontend) npm run lint ;;
      backend)  golangci-lint run ./... ;;
      "")       npm run lint && golangci-lint run ./... ;;
      *)        echo "[forge] invalid scope '{{scope}}'; expected frontend/backend" >&2; exit 1 ;;
    esac

# fmt: auto-format code
fmt scope="":
    #!/usr/bin/env bash
    set -euo pipefail
    case "{{scope}}" in
      frontend) npx prettier --write . ;;
      backend)  gofmt -w . ;;
      "")       npx prettier --write . && gofmt -w . ;;
      *)        echo "[forge] invalid scope '{{scope}}'; expected frontend/backend" >&2; exit 1 ;;
    esac

# check: lint + compile (CI gate)
check scope="":
    #!/usr/bin/env bash
    set -euo pipefail
    case "{{scope}}" in
      frontend) npm run lint && npx tsc --noEmit ;;
      backend)  golangci-lint run ./... ;;
      "")       npm run lint && npx tsc --noEmit && golangci-lint run ./... ;;
      *)        echo "[forge] invalid scope '{{scope}}'; expected frontend/backend" >&2; exit 1 ;;
    esac

# clean: remove build artifacts
clean scope="":
    #!/usr/bin/env bash
    set -euo pipefail
    case "{{scope}}" in
      frontend) rm -rf dist ;;
      backend)  go clean ./... ;;
      "")       rm -rf dist && go clean ./... ;;
      *)        echo "[forge] invalid scope '{{scope}}'; expected frontend/backend" >&2; exit 1 ;;
    esac

# install: install dependencies (idempotent)
install scope="":
    #!/usr/bin/env bash
    set -euo pipefail
    case "{{scope}}" in
      frontend) npm install ;;
      backend)  go mod download ;;
      "")       npm install && go mod download ;;
      *)        echo "[forge] invalid scope '{{scope}}'; expected frontend/backend" >&2; exit 1 ;;
    esac

# ci: full CI pipeline
ci:
    #!/usr/bin/env bash
    set -euo pipefail
    just install
    just compile
    just build
    just test
    just lint

# e2e-setup: install e2e dependencies (idempotent); pass force to always run npm install
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

# e2e-verify: check for unresolved // VERIFY: markers
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
        echo "$matches" >&2
        exit 1
    fi
    echo "OK: no unresolved // VERIFY: markers in tests/e2e/features/{{feature}}/"

# --- end forge standard recipes ---
