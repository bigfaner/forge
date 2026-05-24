# Custom recipes (project-specific, not part of forge standard)

claude-p:
    claude --dangerously-skip-permissions --plugin-dir plugins/forge

# install-forge: build and install forge CLI locally (platform-aware)
install-forge:
    #!/usr/bin/env bash
    set -euo pipefail
    case "$(uname -s)" in
        Linux|Darwin)  bash forge-cli/scripts/install-local.sh ;;
        MINGW*|MSYS*|CYGWIN*) powershell -File forge-cli/scripts/install-local.ps1 ;;
        *) echo "Unsupported OS: $(uname -s)" >&2; exit 1 ;;
    esac

# --- forge standard recipes ---

# compile: type-check and transpile for fast feedback
compile scope="":
    #!/usr/bin/env bash
    set -euo pipefail
    cd forge-cli && go vet ./...

# build: full compile and package
build scope="":
    #!/usr/bin/env bash
    set -euo pipefail
    cd forge-cli && go build ./...

# run: start the service
run scope="":
    #!/usr/bin/env bash
    set -euo pipefail
    cd forge-cli && go run ./cmd/forge

# dev: hot-reload development mode
dev scope="":
    #!/usr/bin/env bash
    set -euo pipefail
    cd forge-cli && go run ./cmd/forge

# unit-test: language-level unit tests
unit-test scope="":
    #!/usr/bin/env bash
    set -euo pipefail
    if command -v gcc &>/dev/null; then
        cd forge-cli && go test -race ./...
    else
        cd forge-cli && CGO_ENABLED=0 go test ./...
    fi

# test: surface-level advanced tests (e2e, integration, etc.)
test journey='':
    #!/usr/bin/env bash
    set -euo pipefail
    feature_flag=""
    if [ -n "{{journey}}" ]; then
        feature_flag="-run TestTC.*$(echo '{{journey}}' | sed 's/.*/\u&/')"
    fi
    cd tests && go test -v -tags=e2e -timeout=10m -json $feature_flag ./... \
      | go-junit-report > results/report.xml 2>/dev/null \
      || go test -v -tags=e2e -timeout=10m $feature_flag ./...

# lint: static analysis
lint scope="":
    #!/usr/bin/env bash
    set -euo pipefail
    cd forge-cli && golangci-lint run ./...

# fmt: auto-format code
fmt scope="":
    #!/usr/bin/env bash
    set -euo pipefail
    cd forge-cli && gofmt -w .

# check: lint + compile (CI gate)
check scope="":
    #!/usr/bin/env bash
    set -euo pipefail
    cd forge-cli && golangci-lint run ./... && go vet ./...

# clean: remove build artifacts
clean scope="":
    #!/usr/bin/env bash
    set -euo pipefail
    cd forge-cli && go clean ./...

# install: install dependencies (idempotent)
install scope="":
    #!/usr/bin/env bash
    set -euo pipefail
    cd forge-cli && go mod download

# ci: full CI pipeline
ci:
    #!/usr/bin/env bash
    set -euo pipefail
    just install
    just compile
    just build
    just unit-test
    just lint

# check-stale-refs: detect old task-cli command references (CI stale reference detection)
check-stale-refs:
    #!/usr/bin/env bash
    set -euo pipefail
    pattern='(^\s*\$?\s*|`)(task (claim|submit|status|query|check-deps|validate-index|verify-task-done|quality-gate|cleanup|feature|prompt|add|index|migrate|validate-specs|record|all-completed|verify-completion|check|validate))\b'
    matches=$(grep -rP "$pattern" plugins/ forge-cli/docs/ --include='*.md' --include='*.json' 2>/dev/null || true)
    if [ -n "$matches" ]; then
        count=$(echo "$matches" | wc -l | tr -d ' ')
        echo "Error: $count stale task-cli reference(s) found" >&2
        echo "$matches" >&2
        exit 1
    fi
    echo "OK: no stale task-cli references"

# test-setup: pre-build forge binary and warm build cache for faster test startup.
# Tests auto-build via TestMain, so this recipe is NOT required before running tests.
# Use it to prime the Go build cache and skip the ~2-5s build during your next test run.
test-setup:
    #!/usr/bin/env bash
    set -euo pipefail
    # Pre-build forge binary for faster test startup (cache optimization)
    cd forge-cli && go build -o bin/forge.exe ./cmd/forge/ && cp bin/forge.exe bin/forge
    # Pre-compile e2e test packages to warm the build cache
    cd tests && go build -tags=e2e ./...
    echo "OK: build cache warmed (optional — tests auto-build via TestMain)"

# probe: check if configured services are healthy
probe path="":
    #!/usr/bin/env bash
    set -euo pipefail
    echo "OK: forge CLI project (no services to probe)"

# test-discover: list all e2e test cases without running them
test-discover:
    #!/usr/bin/env bash
    set -euo pipefail
    cd tests && go test -tags=e2e -list '.*' ./...

# --- end forge standard recipes ---
