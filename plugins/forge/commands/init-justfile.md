---
name: init-justfile
description: Scaffold a Justfile with standard forge targets for the current project.
---

# /init-justfile

Generate a Justfile with standard forge targets as an abstraction layer for test/build commands.

## Prerequisites

**Install just (>= 1.50.0)**

| Platform | Command |
|----------|---------|
| macOS / Linux | `brew install just` |
| Windows (Scoop) | `scoop install just` |
| Windows (winget) | `winget install --id Casey.Just --exact` |
| Cargo (universal) | `cargo install just` |

```bash
just --version  # requires >= 1.50.0 (supports [arg] named option syntax)
```

If version < 1.50.0: `cargo install just`

## Standard Target Contract

| Target | Required | Purpose |
|--------|----------|---------|
| `project-type` | Yes | Return project type identifier (`frontend`/`backend`/`mixed`) |
| `compile` | No | Type-check and transpile for fast feedback |
| `build` | No | Full compile and package |
| `run` | No | Start the service |
| `dev` | No | Hot-reload development mode |
| `test` | Yes | Unit + integration tests |
| `test-e2e` | No | E2E tests |
| `lint` | No | Static analysis |
| `fmt` | No | Auto-format code |
| `check` | No | lint + compile (CI gate) |
| `clean` | No | Remove build artifacts |
| `install` | No | Install dependencies (idempotent) |
| `ci` | No | Full CI pipeline |
| `e2e-setup` | No | Install e2e dependencies (idempotent) |
| `e2e-verify` | No | Check for unresolved `// VERIFY:` markers |

`test-e2e` invocation:

| Invocation | Description |
|------------|-------------|
| `just test-e2e` | Regression tests (`tests/e2e/`) |
| `just test-e2e --feature <slug>` | Specified feature's test scripts |

## Workflow

### Step 1: Detect Project Type

```bash
ls go.mod package.json pyproject.toml setup.py Cargo.toml 2>/dev/null
```

| File | Project Type |
|------|-------------|
| `go.mod` | Go |
| `package.json` | Node.js |
| `pyproject.toml` / `setup.py` | Python |
| `Cargo.toml` | Rust |
| Other | Generic |

### Step 2: Check Existing Files

```bash
ls justfile Justfile Makefile 2>/dev/null
```

- If `justfile` or `Justfile` already exists → ask to overwrite, abort if no
- If `Makefile` exists → read content, migrate existing targets to Justfile

### Step 3: Generate Justfile

Write to `justfile` (lowercase). All templates share the same `test-e2e`; only `test`/`build`/`lint` differ by language.

**test-e2e (all languages):**

```just
# Run e2e tests: "just test-e2e" (regression) or "just test-e2e --feature <slug>" (feature tests)
[arg("feature", long)]
test-e2e feature="":
    #!/usr/bin/env bash
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
```

**e2e-setup (all languages):**

```just
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
```

**e2e-verify (all languages):**

```just
# Requires just >= 1.50.0. [arg("feature", long)] declares a long-form CLI flag:
# --feature <value>. "long" maps the argument to --<name> form; omitting it makes
# the argument positional.
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
```

**Language-specific recipes:**

Go:
```just
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
lint:
    golangci-lint run ./...
fmt:
    gofmt -w .
check:
    #!/usr/bin/env bash
    set -euo pipefail
    golangci-lint run ./...
clean:
    go clean ./...
install:
    go mod download
```

Rust:
```just
project-type:
    @echo "backend"
compile:
    cargo check
build:
    cargo build --release
run:
    cargo run
dev:
    cargo watch -x run
test:
    cargo test
lint:
    cargo clippy -- -D warnings
fmt:
    cargo fmt
check:
    #!/usr/bin/env bash
    set -euo pipefail
    cargo clippy -- -D warnings
clean:
    cargo clean
install:
    cargo fetch
```

Node.js:
```just
project-type:
    @echo "frontend"
compile:
    npx tsc --noEmit
build:
    npm run build
run:
    npm start
dev:
    npm run dev
test:
    npm test
lint:
    npm run lint
fmt:
    npx prettier --write .
check:
    #!/usr/bin/env bash
    set -euo pipefail
    npm run lint && npx tsc --noEmit
clean:
    rm -rf dist
install:
    npm install
```

Python:
```just
project-type:
    @echo "backend"
compile:
    python -m py_compile src/
build:
    python -m build
run:
    python -m src
dev:
    uvicorn src:app --reload
test:
    pytest
lint:
    ruff check .
fmt:
    ruff format .
check:
    #!/usr/bin/env bash
    set -euo pipefail
    ruff check .
clean:
    rm -rf build/ dist/ *.egg-info
install:
    pip install -e .
```

Generic:
```just
project-type:
    @echo "backend"
compile:
    echo "TODO: implement compile recipe"
build:
    echo "TODO: implement build recipe"
run:
    echo "TODO: implement run recipe"
dev:
    echo "TODO: implement dev recipe"
test:
    echo "TODO: implement test recipe"
lint:
    echo "TODO: implement lint recipe"
fmt:
    echo "TODO: implement fmt recipe"
check:
    echo "TODO: implement check recipe"
clean:
    echo "TODO: implement clean recipe"
install:
    echo "TODO: implement install recipe"
```

### Step 4: Output Confirmation

```
Created justfile with standard forge targets (Go project)

Targets:
  just project-type               → echo "backend"
  just compile                    → go vet ./...
  just build                      → go build ./...
  just run                        → go run .
  just dev                        → go run . --dev
  just test                       → go test -race ./...
  just test-e2e                   → regression tests in tests/e2e/
  just test-e2e --feature <slug>  → feature tests in tests/e2e/<slug>/
  just lint                       → golangci-lint run ./...
  just fmt                        → gofmt -w .
  just check                      → golangci-lint run ./...
  just clean                      → go clean ./...
  just install                    → go mod download
  just ci                         → install + compile + build + test + lint
  just e2e-setup                  → install e2e deps (idempotent)
  just e2e-verify --feature <slug> → check for unresolved // VERIFY: markers

Edit justfile to customize commands for your project.
task all-completed will now use `just test` automatically.
```

## Recipe Templates

Each project type (backend, frontend, mixed) has its own set of 15 recipe string literals.
These templates are assembled by the project-type detection logic (Step 1) to produce the final justfile.

### Backend Template (pure backend projects: Go/Rust/Python)

No scope parameters. Each recipe directly calls the Go toolchain command.

```just
# --- forge standard recipes ---

# project-type: return project type identifier
project-type:
    @echo "backend"

# compile: type-check and transpile for fast feedback
compile:
    go vet ./...

# build: full compile and package
build:
    go build ./...

# run: start the service
run:
    go run .

# dev: hot-reload development mode
dev:
    go run . --dev

# test: unit + integration tests
test:
    go test -race ./...

# test-e2e: end-to-end tests
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

# lint: static analysis
lint:
    golangci-lint run ./...

# fmt: auto-format code
fmt:
    gofmt -w .

# check: lint + compile (CI gate)
check:
    #!/usr/bin/env bash
    set -euo pipefail
    golangci-lint run ./...

# clean: remove build artifacts
clean:
    go clean ./...

# install: install dependencies (idempotent)
install:
    go mod download

# ci: full CI pipeline
ci:
    #!/usr/bin/env bash
    set -euo pipefail
    just install
    just compile
    just build
    just test
    just lint

# e2e-setup: install e2e dependencies (idempotent)
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

# e2e-verify: check for unresolved // VERIFY: markers
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
```

## Notes

- **just >= 1.50.0**: `[arg("feature", long)]` generates `--feature <value>` named option, must pass a value when invoked
- Callers (CI, `task all-completed`) are responsible for passing the slug: `just test-e2e --feature <slug>`
- When migrating from Makefile, preserve original command logic, only adjust format
