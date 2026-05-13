# Custom recipes (project-specific, not part of forge standard)

claude:
    claude --dangerously-skip-permissions

claude-c:
    claude --dangerously-skip-permissions -c

claude-w name="":
    claude --dangerously-skip-permissions -w "{{name}}"

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

# project-type: return project type identifier
project-type:
    @echo "backend"

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
    cd forge-cli && go run .

# dev: hot-reload development mode
dev scope="":
    #!/usr/bin/env bash
    set -euo pipefail
    cd forge-cli && go run .

# test: unit + integration tests
test scope="":
    #!/usr/bin/env bash
    set -euo pipefail
    cd forge-cli && go test -race ./...

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
    # Match standalone `task <subcommand>` used as CLI invocation (e.g. in shell snippets)
    # but exclude: `forge task <subcommand>` (correct), natural language ("task status"),
    # template variables ("task index"), and markdown table decorations.
    pattern='(^\s*\$?\s*|`)(task (claim|submit|status|query|check-deps|validate-index|verify-task-done|quality-gate|cleanup|feature|prompt|add|index|migrate|validate-specs|record|all-completed|verify-completion|check|validate))\b'
    matches=$(grep -rP "$pattern" plugins/ forge-cli/docs/ --include='*.md' --include='*.json' 2>/dev/null || true)
    if [ -n "$matches" ]; then
        count=$(echo "$matches" | wc -l | tr -d ' ')
        echo "Error: $count stale task-cli reference(s) found" >&2
        echo "$matches" >&2
        exit 1
    fi
    echo "OK: no stale task-cli references"

# --- end forge standard recipes ---
