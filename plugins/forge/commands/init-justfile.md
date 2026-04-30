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

Check for project marker files and classify the project type:

```bash
ls package.json go.mod Cargo.toml pyproject.toml 2>/dev/null
```

**Detection signal mapping:**

| Marker File | Signal |
|-------------|--------|
| `package.json` | frontend |
| `go.mod` | backend |
| `Cargo.toml` | backend |
| `pyproject.toml` | backend |

**Classification algorithm:**

1. Check for each marker file's existence in the project root.
2. Collect detected signals into two sets: `frontend_signals` and `backend_signals`.
3. Classify:
   - Exactly one frontend signal AND exactly one backend signal → **`mixed`**
   - Exactly one frontend signal, no backend signals → **`frontend`**
   - Exactly one backend signal, no frontend signals → **`backend`**
   - Neither set has signals → **Error**: "Error: no known project markers detected (expected one of: package.json, go.mod, Cargo.toml, pyproject.toml)" — abort, do NOT generate a justfile.
   - Multiple frontend signals (e.g. `package.json` + another) → **Error**: "Error: multiple frontend markers detected — not supported" — abort.
   - Multiple backend signals (e.g. `go.mod` + `Cargo.toml`) → **Error**: "Error: multiple backend markers detected — not supported" — abort.

Supported configurations: backend only, frontend only, one backend + one frontend (mixed).

**For `mixed` projects, also detect root paths:**

```bash
# Find frontend root (package.json, excluding node_modules)
find . -name package.json -not -path '*/node_modules/*' -maxdepth 3 | head -1 | xargs dirname

# Find backend root (go.mod / Cargo.toml / pyproject.toml)
find . \( -name go.mod -o -name Cargo.toml -o -name pyproject.toml \) -maxdepth 3 | head -1 | xargs dirname
```

Record these as `FRONTEND_DIR` and `BACKEND_DIR` (e.g. `./frontend`, `./backend`). Use `.` if the marker is in the project root.

### Step 2: Check Existing Justfile

```bash
ls justfile Justfile 2>/dev/null
```

- If `justfile` or `Justfile` already exists:
  - Check for boundary markers (`# --- forge standard recipes ---` / `# --- end forge standard recipes ---`).
  - **If boundary markers exist**: proceed to Step 3 (boundary marker merge). No confirmation needed — only the marked section will be replaced; custom recipes outside markers are preserved.
  - **If boundary markers do NOT exist** (user's justfile has no forge markers):
    - If `--force` flag was provided: skip confirmation, proceed to Step 3 (will add boundary markers and replace entire content).
    - If `--force` flag was NOT provided: prompt the user: "A justfile already exists without forge markers. Overwrite? (y/n)". If user declines, abort without modifying the file.
- If no justfile exists: proceed to Step 3 (create new file).

**`--force` flag**: Use `/init-justfile --force` to skip all interactive confirmation prompts. This is required for agent workflows where no human is present to respond to prompts. When `--force` is active, the command runs non-interactively: it merges within boundary markers if they exist, or overwrites the entire file if they don't.

### Step 3: Assemble and Write Justfile

Select the template from `plugins/forge/references/justfile-templates/`:

| Marker file | Template |
|-------------|----------|
| `go.mod` | `go.just` |
| `Cargo.toml` | `rust.just` |
| `pyproject.toml` | `python.just` |
| `package.json` only | `node.just` |
| mixed | `mixed.just` |
| none matched | `generic.just` |

Write to `justfile` (lowercase).

**Boundary marker merge**: When an existing justfile contains `# --- forge standard recipes ---` / `# --- end forge standard recipes ---` markers:
1. Read the existing justfile content.
2. Replace everything between (and including) the markers with the new template content.
3. Preserve all content outside the markers (user custom recipes).
4. Write the merged result back to `justfile`.

**New justfile**: Write the selected template as the new file.

**For `mixed` projects**: after copying `mixed.just`, replace the two placeholder variables with the detected paths:

- `FRONTEND_DIR` → detected frontend root (e.g. `./frontend`)
- `BACKEND_DIR` → detected backend root (e.g. `./backend`)

### Step 4: Verify Generated Commands

Run the following to confirm the justfile is valid and recipes are callable:

```bash
just --list
just project-type
just --dry-run compile
```

If any command fails, report the error and do not claim success.

### Step 5: Output Confirmation

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

## Notes

- **just >= 1.50.0**: `[arg("feature", long)]` generates `--feature <value>` named option, must pass a value when invoked
- Callers (CI, `task all-completed`) are responsible for passing the slug: `just test-e2e --feature <slug>`
- When migrating from Makefile, preserve original command logic, only adjust format
- **e2e tests use `npx tsx` regardless of project language**: forge e2e test scripts are always written in TypeScript and run via `tsx`. This is intentional — the e2e layer is language-agnostic. Backend projects (Go/Rust/Python) still need Node.js available in the environment for `just test-e2e` and `just e2e-setup`.
