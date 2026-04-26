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
| `test` | Yes | Unit + integration tests |
| `test-e2e` | No | E2E tests |
| `build` | No | Compile/package |
| `lint` | No | Static analysis |

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
        scripts_dir="docs/features/{{feature}}/testing/scripts"
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

**Language-specific recipes:**

Go:
```just
test:
    go test -race ./...
build:
    go build ./...
lint:
    golangci-lint run ./...
```

Rust:
```just
test:
    cargo test
build:
    cargo build --release
lint:
    cargo clippy -- -D warnings
```

Node.js:
```just
test:
    npm test
build:
    npm run build
lint:
    npm run lint
```

Python:
```just
test:
    pytest
build:
    python -m build
lint:
    ruff check .
```

Generic:
```just
test:
    echo "TODO: implement test recipe"
build:
    echo "TODO: implement build recipe"
lint:
    echo "TODO: implement lint recipe"
```

### Step 4: Output Confirmation

```
Created justfile with standard forge targets (Go project)

Targets:
  just test                       → go test -race ./...
  just test-e2e                   → regression tests in tests/e2e/
  just test-e2e --feature <slug>  → feature tests in docs/features/<slug>/testing/scripts/
  just build                      → go build ./...
  just lint                       → golangci-lint run ./...

Edit justfile to customize commands for your project.
task all-completed will now use `just test` automatically.
```

## Notes

- **just >= 1.50.0**: `[arg("feature", long)]` generates `--feature <value>` named option, must pass a value when invoked
- Callers (CI, `task all-completed`) are responsible for passing the slug: `just test-e2e --feature <slug>`
- When migrating from Makefile, preserve original command logic, only adjust format
