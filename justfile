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

# test: unit + integration tests
test scope="":
    #!/usr/bin/env bash
    set -euo pipefail
    if command -v gcc &>/dev/null; then
        cd forge-cli && go test -race ./...
    else
        cd forge-cli && CGO_ENABLED=0 go test ./...
    fi

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
    just test
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

# test-e2e: end-to-end tests (go-test profile)
[arg("feature", long)]
test-e2e feature="":
    #!/usr/bin/env bash
    set -euo pipefail
    feature_flag=""
    if [ -n "{{feature}}" ]; then
        feature_flag="-run TestTC.*$(echo '{{feature}}' | sed 's/.*/\u&/')"
    fi
    cd tests/e2e && go test -v -tags=e2e -timeout=10m -json $feature_flag \
      | go-junit-report > results/report.xml 2>/dev/null \
      || go test -v -tags=e2e -timeout=10m $feature_flag

# e2e-setup: verify compilation (idempotent, go-test profile)
e2e-setup force="":
    #!/usr/bin/env bash
    set -euo pipefail
    cd forge-cli && go build -o bin/forge.exe ./cmd/forge/ && cp bin/forge.exe bin/forge
    cd tests/e2e && go build -tags=e2e ./...
    echo "OK: compilation verified"

# e2e-verify: check for unresolved // VERIFY: markers (go-test profile)
[arg("feature", long)]
e2e-verify feature="":
    #!/usr/bin/env bash
    set -euo pipefail
    if [ -z "{{feature}}" ]; then
        echo "Usage: just e2e-verify --feature <slug>" >&2
        exit 1
    fi
    search_dir="tests/e2e/features/{{feature}}"
    if [ ! -d "$search_dir" ]; then
        search_dir="tests/e2e"
    fi
    matches=$(grep -rn '// VERIFY:' "$search_dir/" --include='*_test.go' || true)
    if [ -n "$matches" ]; then
        count=$(echo "$matches" | wc -l | tr -d ' ')
        echo "Error: $count unresolved // VERIFY: marker(s) in $search_dir/" >&2
        echo "$matches" >&2
        exit 1
    fi
    echo "OK: no unresolved // VERIFY: markers in $search_dir/"

# e2e-compile: compile-check e2e test files
e2e-compile:
    #!/usr/bin/env bash
    set -euo pipefail
    cd tests/e2e && go build -tags=e2e ./...
    echo "OK: Go compilation passed"

# e2e-discover: list all e2e test cases without running them
e2e-discover:
    #!/usr/bin/env bash
    set -euo pipefail
    cd tests/e2e && go test -tags=e2e -list '.*' ./...

# --- end forge standard recipes ---
