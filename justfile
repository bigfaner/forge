# Custom recipes (project-specific, not part of forge standard)

claude:
    claude --dangerously-skip-permissions

claude-c:
    claude --dangerously-skip-permissions -c

claude-w name="":
    claude --dangerously-skip-permissions -w "{{name}}"

claude-p:
    claude --dangerously-skip-permissions --plugin-dir plugins/forge

# install-task: build and install task CLI locally (platform-aware)
install-task:
    #!/usr/bin/env bash
    set -euo pipefail
    case "$(uname -s)" in
        Linux|Darwin)  bash task-cli/scripts/install-local.sh ;;
        MINGW*|MSYS*|CYGWIN*) powershell -File task-cli/scripts/install-local.ps1 ;;
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
    cd task-cli && go vet ./...

# build: full compile and package
build scope="":
    #!/usr/bin/env bash
    set -euo pipefail
    cd task-cli && go build ./...

# run: start the service
run scope="":
    #!/usr/bin/env bash
    set -euo pipefail
    cd task-cli && go run .

# dev: hot-reload development mode
dev scope="":
    #!/usr/bin/env bash
    set -euo pipefail
    cd task-cli && go run .

# test: unit + integration tests
test scope="":
    #!/usr/bin/env bash
    set -euo pipefail
    cd task-cli && go test -race ./...

# test-e2e: end-to-end tests (profile-aware)
[arg("feature", long)]
test-e2e feature="":
    #!/usr/bin/env bash
    set -euo pipefail
    profile=$(grep -A1 'test-profiles:' .forge/config.yaml 2>/dev/null | grep -oE '\w[-\w]*' | head -1 || echo "web-playwright")
    case "$profile" in
      web-playwright)
        if [ -z "{{feature}}" ]; then
            [ ! -d tests/e2e/node_modules ] && npm install --prefix tests/e2e
            cd tests/e2e && npx playwright test
        else
            feature_config="tests/e2e/features/{{feature}}/playwright.config.ts"
            if [ -f "$feature_config" ]; then
                cd tests/e2e/features/{{feature}} && npx playwright test --config=playwright.config.ts
            else
                cd tests/e2e && E2E_FEATURE=1 npx playwright test features/{{feature}}/
            fi
        fi
        ;;
      *)
        echo "Error: unknown test profile '$profile' in .forge/config.yaml" >&2
        echo "Supported profiles: web-playwright, go-test, maestro, java-junit, rust-test, pytest" >&2
        exit 1 ;;
    esac

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
    cd task-cli && golangci-lint run ./...

# fmt: auto-format code
fmt scope="":
    #!/usr/bin/env bash
    set -euo pipefail
    cd task-cli && gofmt -w .

# check: lint + compile (CI gate)
check scope="":
    #!/usr/bin/env bash
    set -euo pipefail
    cd task-cli && golangci-lint run ./... && go vet ./...

# clean: remove build artifacts
clean scope="":
    #!/usr/bin/env bash
    set -euo pipefail
    cd task-cli && go clean ./...

# install: install dependencies (idempotent)
install scope="":
    #!/usr/bin/env bash
    set -euo pipefail
    cd task-cli && go mod download

# ci: full CI pipeline
ci:
    #!/usr/bin/env bash
    set -euo pipefail
    just install
    just compile
    just build
    just test
    just lint

# e2e-setup: install e2e dependencies (idempotent, profile-aware)
e2e-setup force="":
    #!/usr/bin/env bash
    set -euo pipefail
    profile=$(grep -A1 'test-profiles:' .forge/config.yaml 2>/dev/null | grep -oE '\w[-\w]*' | head -1 || echo "web-playwright")
    case "$profile" in
      web-playwright)
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
        echo "OK: e2e dependencies ready (web-playwright)"
        ;;
      *)
        echo "Error: unknown test profile '$profile' in .forge/config.yaml" >&2
        exit 1 ;;
    esac

# e2e-verify: check for unresolved // VERIFY: markers (profile-aware file extension)
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
    ext=".spec.ts"
    profile=$(grep -A1 'test-profiles:' .forge/config.yaml 2>/dev/null | grep -oE '\w[-\w]*' | head -1 || echo "web-playwright")
    case "$profile" in
      web-playwright) ext=".spec.ts" ;;
      go-test) ext="_test.go" ;;
      java-junit) ext="Test.java" ;;
      rust-test) ext=".rs" ;;
      pytest) ext=".py" ;;
      maestro) ext=".yaml" ;;
    esac
    matches=$(grep -rn '// VERIFY:' "tests/e2e/features/{{feature}}/" --include="*${ext}" || true)
    if [ -n "$matches" ]; then
        count=$(echo "$matches" | wc -l | tr -d ' ')
        echo "Error: $count unresolved // VERIFY: marker(s) in tests/e2e/features/{{feature}}/" >&2
        echo "$matches" >&2
        exit 1
    fi
    echo "OK: no unresolved // VERIFY: markers in tests/e2e/features/{{feature}}/"

# e2e-compile: compile-check e2e test files (profile-aware)
e2e-compile:
    #!/usr/bin/env bash
    set -euo pipefail
    profile=$(grep -A1 'test-profiles:' .forge/config.yaml 2>/dev/null | grep -oE '\w[-\w]*' | head -1 || echo "web-playwright")
    case "$profile" in
      web-playwright) cd tests/e2e && npx tsc --noEmit && echo "OK: TypeScript compilation passed" ;;
      go-test) go build ./tests/e2e/... && echo "OK: Go compilation passed" ;;
      maestro)
        errors=0
        for f in tests/e2e/features/**/*.yaml tests/e2e/*.yaml; do
          [ -f "$f" ] || continue
          if ! maestro test --dry-run "$f" 2>/dev/null; then
            echo "Error: invalid maestro flow: $f" >&2
            errors=$((errors+1))
          fi
        done
        [ "$errors" -eq 0 ] && echo "OK: all maestro flows valid" || exit 1 ;;
      java-junit) mvn test-compile -pl tests/e2e && echo "OK: Java compilation passed" ;;
      rust-test) cargo build --test e2e && echo "OK: Rust compilation passed" ;;
      pytest) python -m compileall tests/e2e/ -q && echo "OK: Python compilation passed" ;;
      *)
        echo "Error: unknown test profile '$profile' in .forge/config.yaml" >&2
        exit 1 ;;
    esac

# e2e-discover: list all e2e test cases without running them (profile-aware)
e2e-discover:
    #!/usr/bin/env bash
    set -euo pipefail
    profile=$(grep -A1 'test-profiles:' .forge/config.yaml 2>/dev/null | grep -oE '\w[-\w]*' | head -1 || echo "web-playwright")
    case "$profile" in
      web-playwright) cd tests/e2e && npx playwright test --list ;;
      go-test) go test ./tests/e2e/... -list '.*' -tags=e2e ;;
      maestro) find tests/e2e -name '*.yaml' -not -name 'config.yaml' | sort ;;
      java-junit) mvn test -pl tests/e2e -DdryRun=true ;;
      rust-test) cargo test --test e2e -- --list ;;
      pytest) python -m pytest tests/e2e/ --collect-only -q ;;
      *)
        echo "Error: unknown test profile '$profile' in .forge/config.yaml" >&2
        exit 1 ;;
    esac

# --- end forge standard recipes ---
