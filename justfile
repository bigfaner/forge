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

# test-discover: list all e2e test cases without running them
test-discover:
    #!/usr/bin/env bash
    set -euo pipefail
    cd tests && go test -tags=cli_functional -list '.*' ./...

# --- forge standard recipes ---

# compile: type-check for fast feedback
# user-customized
[group("go")]
compile:
    #!/usr/bin/env bash
    set -euo pipefail
    cd forge-cli && go vet ./...

# build: full compile and package
# user-customized
[group("go")]
build:
    #!/usr/bin/env bash
    set -euo pipefail
    cd forge-cli && go build ./...

# fmt: auto-format code
# user-customized
[group("go")]
fmt:
    #!/usr/bin/env bash
    set -euo pipefail
    cd forge-cli && gofmt -w .

# install: install dependencies (idempotent)
# user-customized
[group("go")]
install:
    #!/usr/bin/env bash
    set -euo pipefail
    cd forge-cli && go mod download

# clean: remove build artifacts
# user-customized
[group("go")]
clean:
    #!/usr/bin/env bash
    set -euo pipefail
    cd forge-cli && go clean ./...

# ci: full CI pipeline
# user-customized
[group("go")]
ci:
    #!/usr/bin/env bash
    set -euo pipefail
    just install
    just compile
    just build
    just unit-test
    just lint

# lint: static analysis
# user-customized
[group("go-test")]
lint:
    #!/usr/bin/env bash
    set -euo pipefail
    cd forge-cli && golangci-lint run ./...

# check: lint + compile (CI gate)
# user-customized
[group("go-test")]
check:
    #!/usr/bin/env bash
    set -euo pipefail
    cd forge-cli && golangci-lint run ./... && go vet ./...

# unit-test: language-level unit tests
# user-customized
[group("go-test")]
unit-test:
    #!/usr/bin/env bash
    set -euo pipefail
    if command -v gcc &>/dev/null; then
        cd forge-cli && go test -race ./...
    else
        cd forge-cli && CGO_ENABLED=0 go test ./...
    fi

# test: surface-level CLI functional tests
# user-customized
[group("cli")]
test journey='':
    #!/usr/bin/env bash
    set -euo pipefail
    feature_flag=""
    if [ -n "{{journey}}" ]; then
        feature_flag="-run TestTC.*$(echo '{{journey}}' | sed 's/.*/\u&/')"
    fi
    cd tests && go test -v -tags=cli_functional -timeout=10m -json $feature_flag ./... \
      | go-junit-report > results/report.xml 2>/dev/null \
      || go test -v -tags=cli_functional -timeout=10m $feature_flag ./...

# teardown: cleanup after tests (no-op for CLI surface)
# user-customized
[group("cli")]
teardown:
    #!/usr/bin/env bash
    echo "OK: forge CLI project (no teardown needed)"

# --- end forge standard recipes ---
