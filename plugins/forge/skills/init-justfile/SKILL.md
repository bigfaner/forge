---
name: init-justfile
description: Scaffold a Justfile with standard forge targets for the current project.
allowed_tools: ["Bash", "Read", "Write", "Edit"]
disable-model-invocation: true
argument-hints: "[--lang go|rust|python|node] [--type frontend|backend|mixed] [--force]"
---

# /init-justfile

MANUAL-ONLY. Do NOT auto-invoke — only when user explicitly asks to `/init-justfile`.

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

## Parameters

| Parameter | Values | Default | Description |
|-----------|--------|---------|-------------|
| `--lang` | `go`, `rust`, `python`, `node` | (auto-detect) | Override language detection |
| `--type` | `frontend`, `backend`, `mixed` | (auto-detect) | Override project type |
| `--force` | (flag) | false | Overwrite existing justfile without confirmation |

## Standard Target Contract

| Target | Required | Purpose |
|--------|----------|---------|
| `project-type` | Yes | Return project type identifier (`frontend`/`backend`/`mixed`) |
| `compile` | No | Type-check and transpile for fast feedback |
| `build` | No | Full compile and package |
| `run` | No | Start the service |
| `dev` | No | Hot-reload development mode |
| `test` | Yes | Unit + integration tests |
| `test-e2e` | No | E2E tests; `--feature <slug>` for single feature |
| `lint` | No | Static analysis |
| `fmt` | No | Auto-format code |
| `check` | No | lint + compile (CI gate) |
| `clean` | No | Remove build artifacts |
| `install` | No | Install dependencies (idempotent) |
| `ci` | No | Full CI pipeline |
| `e2e-setup` | No | Install e2e dependencies (idempotent) |
| `e2e-verify` | No | Check for unresolved `// VERIFY:` markers |

## Workflow

```
0. Resolve test profile → 1. Detect project type + language + entry points → 2. Check existing justfile → 3. Assemble and write → 4. Verify and self-correct → 5. Output confirmation
```

### Step 0: Resolve Test Profile

1. **Resolve profile**: Run `task profile` to get the active test profile(s). This reads `.forge/config.yaml`, falls back to project structure detection.
2. **On failure** (output shows `PROFILE: (none)`): ask the user to choose from known profiles (`web-playwright`, `go-test`, `maestro`, `java-junit`, `rust-test`, `pytest`). Run `task profile set <name>` to persist their choice.
3. **Load profile manifest**: Read `plugins/forge/profiles/<profile-name>/manifest.yaml`.
4. **Load justfile recipes**: Read `plugins/forge/profiles/<profile-name>/justfile-recipes` for the profile-specific e2e recipe bodies.

The `test-e2e`, `e2e-setup`, and `e2e-verify` recipes are generated from the profile's `justfile-recipes` file, not from the language template.

<HARD-RULE>
Do NOT silently default to any profile. If `task profile` returns no result and the user cannot decide, abort the skill.
</HARD-RULE>

### Step 1: Detect Project Type, Language, and Entry Points

```bash
just --version 2>/dev/null | awk '{print $2}' | awk -F. '$1 > 1 || ($1 == 1 && $2 >= 50)' | grep -q .
# If exit code != 0: output "Error: just >= 1.50.0 required — run `cargo install just`" and abort.
```

#### Parameter override

If `--lang` and/or `--type` are provided, skip all or part of auto-detection:

1. **`--lang` and `--type` both set**: Skip all detection (1a/1b/1c). Use `--lang` for template selection, `--type` for project-type output. Use placeholder defaults from the table below.
2. **`--lang` only**: Skip 1a (marker detection) and 1b/1c (language-specific detection). Run a simplified 1a just to determine `--type` if no markers exist, otherwise use detected type.
3. **`--type` only**: Run 1a to detect language from markers, but override the project-type classification with `--type`. Run 1b/1c for the detected language.
4. **Neither set**: Full auto-detection (existing behavior).

**Placeholder defaults when skipping detection:**

| Placeholder | Go | Rust | Python | Node |
|-------------|-----|------|--------|------|
| `ENTRY_POINT` | `.` | `` (empty) | `main.py` | N/A |
| `DEV_COMMAND` | N/A | N/A | `python main.py` | N/A |
| `ENTRY_SCRIPT` | N/A | N/A | N/A | `dev` |

For `mixed` projects via parameters, `FRONTEND_DIR` and `BACKEND_DIR` both default to `.`.

#### 1a. Project type detection

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
2. Count frontend vs backend signals and classify:
   - Exactly one frontend signal AND exactly one backend signal → **`mixed`**
   - Exactly one frontend signal, no backend signals → **`frontend`**
   - Exactly one backend signal, no frontend signals → **`backend`**
   - Neither set has signals → **Error**: "Error: no known project markers detected (expected one of: package.json, go.mod, Cargo.toml, pyproject.toml)" — abort, do NOT generate a justfile.
   - Multiple backend signals (e.g. `go.mod` + `Cargo.toml`) → **Error**: "Error: multiple backend markers detected — not supported" — abort.

**For `mixed` projects, also detect root paths:**

```bash
find . -name package.json -not -path '*/node_modules/*' -maxdepth 3 | head -1 | xargs dirname
find . \( -name go.mod -o -name Cargo.toml -o -name pyproject.toml \) -maxdepth 3 | head -1 | xargs dirname
```

Record these as `FRONTEND_DIR` and `BACKEND_DIR` (e.g. `./frontend`, `./backend`). Use `.` if the marker is in the project root.

#### 1b. Backend language + entry point detection

For **backend** and **mixed** projects, detect the backend language (already known from marker file) and the entry point:

| Language | Marker | Entry point detection (`BACKEND_ENTRY`) | `ENTRY_POINT` placeholder value | `DEV_COMMAND` placeholder value |
|----------|--------|---------------------------------------|-------------------------------|-------------------------------|
| Go | `go.mod` | `ls cmd/*/main.go` → `cmd/<name>/main.go`; else `ls main.go` → `.` | `cmd/server/main.go` or `.` | N/A (uses `BACKEND_DEV` row) |
| Rust | `Cargo.toml` | `grep '\[\[bin\]\]' Cargo.toml` → `--bin <name>`; else empty | `--bin server` or `` (empty) | N/A (uses `BACKEND_DEV` row) |
| Python | `pyproject.toml` | `ls src/__init__.py` → `-m src`; `ls main.py` → `main.py`; `ls app.py` → `app.py` | `-m src` / `main.py` / `app.py` | `uvicorn src:app --reload` or `python -m src --reload` |

Record `BACKEND_ENTRY` from the detected entry point. For Python `DEV_COMMAND`: use `uvicorn src:app --reload` if uvicorn is available, else `python -m src --reload`.

#### 1c. Frontend run script detection

For **frontend** and **mixed** projects, detect available npm scripts:

```bash
node -e "const s=JSON.parse(require('fs').readFileSync('package.json')).scripts||{}; console.log(s.start ? 'start' : s.preview ? 'preview' : 'dev')"
```

Record as `FRONTEND_RUN_SCRIPT` (e.g. `start`, `preview`, or `dev`).

### Step 2: Check Existing Justfile

```bash
ls justfile Justfile 2>/dev/null
```

- If `justfile` or `Justfile` already exists:
  - Check for boundary markers (`# --- forge standard recipes ---` / `# --- end forge standard recipes ---`).
  - **If boundary markers exist**: proceed to Step 3 (boundary marker merge). No confirmation needed — only the marked section will be replaced; custom recipes outside markers are preserved.
  - **If boundary markers do NOT exist** (user's justfile has no forge markers):
    - If `--force` flag was provided: skip confirmation, proceed to Step 3 (merge within boundary markers if they exist, or overwrite entire file if they don't).
    - If `--force` flag was NOT provided: prompt the user: "A justfile already exists without forge markers. Overwrite? (y/n)". If user declines, abort without modifying the file.
- If no justfile exists: proceed to Step 3 (create new file).

### Step 3: Assemble and Write Justfile

#### Template selection

If `--lang` is provided, select template directly:

| `--lang` value | Template |
|----------------|----------|
| `go` | `plugins/forge/skills/init-justfile/templates/go.just` |
| `rust` | `plugins/forge/skills/init-justfile/templates/rust.just` |
| `python` | `plugins/forge/skills/init-justfile/templates/python.just` |
| `node` | `plugins/forge/skills/init-justfile/templates/node.just` |
| (mixed via `--type mixed`) | `plugins/forge/skills/init-justfile/templates/mixed.just` |

If `--lang` is not provided, detect from marker files:

| Marker file | Template |
|-------------|----------|
| `go.mod` | `plugins/forge/skills/init-justfile/templates/go.just` |
| `Cargo.toml` | `plugins/forge/skills/init-justfile/templates/rust.just` |
| `pyproject.toml` | `plugins/forge/skills/init-justfile/templates/python.just` |
| `package.json` only | `plugins/forge/skills/init-justfile/templates/node.just` |
| mixed | `plugins/forge/skills/init-justfile/templates/mixed.just` |
| none matched | `plugins/forge/skills/init-justfile/templates/generic.just` |

Write to `justfile` (lowercase).

#### Placeholder substitution

For **single-language templates** (`go.just`, `rust.just`, `python.just`, `node.just`):

| Placeholder | Scope | Replaced with | Example |
|-------------|-------|--------------|---------|
| `ENTRY_POINT` | Go/Rust | `BACKEND_ENTRY` from Step 1b | `cmd/server/main.go`, `.`, `--bin server`, `` (empty) |
| `DEV_COMMAND` | Python only | Dev server command from Step 1b | `uvicorn src:app --reload`, `python -m src --reload` |
| `ENTRY_SCRIPT` | Node only | `FRONTEND_RUN_SCRIPT` from Step 1c | `start` / `preview` / `dev` |

For **mixed projects** (`mixed.just`):

Replace `FRONTEND_DIR` and `BACKEND_DIR` with the paths detected in Step 1a. Additionally, replace all `BACKEND_*` and `FRONTEND_*` placeholders based on detected backend language:

| Placeholder | Go | Rust | Python |
|-------------|-----|------|--------|
| `BACKEND_COMPILE` | `go vet ./...` | `cargo check` | `python -m py_compile src/` |
| `BACKEND_BUILD` | `go build ./...` | `cargo build --release` | `python -m build` |
| `BACKEND_RUN` | `go run <BACKEND_ENTRY>` | `cargo run <BACKEND_ENTRY>` | `python <BACKEND_ENTRY>` |
| `BACKEND_DEV` | `go run <BACKEND_ENTRY>` | `cargo run <BACKEND_ENTRY>` | `uvicorn src:app --reload` |
| `BACKEND_TEST` | `go test ./...` | `cargo test` | `pytest` |
| `BACKEND_LINT` | `golangci-lint run ./...` | `cargo clippy -- -D warnings` | `ruff check .` |
| `BACKEND_FMT` | `gofmt -w .` | `cargo fmt` | `ruff format .` |
| `BACKEND_CLEAN` | `go clean ./...` | `cargo clean` | `rm -rf build/ dist/ *.egg-info` |
| `BACKEND_INSTALL` | `go mod download` | `cargo fetch` | `pip install -e .` |
| `FRONTEND_RUN` | `npm run <FRONTEND_RUN_SCRIPT>` (value from Step 1c) | (same) | (same) |
| `FRONTEND_DEV` | `npm run dev` (all backend langs) | (same) | (same) |

Note: `go test` omits `-race` flag by default — it requires CGO which is unavailable on some platforms (notably Windows). See Step 4c for auto-detection.

#### Boundary marker merge

When markers exist (`# --- forge standard recipes ---` / `# --- end forge standard recipes ---`), replace everything between them (inclusive) with the new template, preserving user recipes outside. Otherwise write the full template as a new file.

### Step 4: Verify and Self-Correct

Two-phase verification: `--dry-run` catches syntax/structure errors, actual execution catches runtime errors (missing scripts, wrong entry points, unavailable tools).

#### 4a. Phase 1 — Dry-run (syntax check)

Run each recipe with `--dry-run` to verify recipe syntax, variable expansion, and command structure:

```bash
just --list
just --dry-run compile
just --dry-run build
just --dry-run test
just --dry-run run
just --dry-run dev
just --dry-run install
just --dry-run lint
just --dry-run fmt
just --dry-run check
just --dry-run e2e-setup
```

For mixed projects, also verify `--dry-run compile` and `--dry-run run` with `backend`/`frontend` scope arguments.

Fix any syntax failures before proceeding to Phase 2.

#### 4b. Phase 2 — Actual execution (runtime check)

Execute each recipe for real to catch runtime errors. Recipes are classified by execution safety:

| Category | Recipes | Method |
|----------|---------|--------|
| **Safe** (fast, no side effects) | `project-type`, `compile`, `lint`, `check` | Execute directly |
| **Destructive** (modifies files or creates artifacts) | `build`, `fmt`, `clean` | Execute directly (artifacts can be cleaned; fmt changes are welcome) |
| **Idempotent** (installs dependencies) | `install`, `e2e-setup` | Execute directly |
| **Long-running** (starts servers) | `run`, `dev` | Execute with timeout (10s), kill after timeout — success = process still alive at timeout |
| **Expensive** (runs full test suite) | `test`, `test-e2e` | Skip actual execution; verified by `--dry-run` only |

For long-running recipes (`run`, `dev`): execute via `timeout 10 just <recipe> 2>&1 || true`. A crash before timeout ("missing script", "can't load package") is a runtime failure. For mixed, also verify scoped variants.

#### 4c. Self-correction rules

When a recipe fails in Phase 2, analyze the error and apply corrections:

| Error Pattern | Recipe | Fix |
|---------------|--------|-----|
| `npm error Missing script: "start"` | `run` (node/mixed) | Replace `npm run start` → `npm run preview` in justfile, retry |
| `npm error Missing script: "preview"` | `run` (node/mixed) | Replace → `npm run dev` in justfile, retry |
| `npm error Missing script: "dev"` | `dev` (node/mixed) | Replace → `npm run start` in justfile, retry |
| `can't load package: no Go files` | `run`/`dev`/`compile` (go) | Scan for `cmd/*/main.go`, update entry point in justfile, retry |
| `CGO_ENABLED=1` available | `test` (go) | Add `-race` flag to `go test` recipe for race detection, retry |
| `command not found: golangci-lint` | `lint`/`check` (go) | In `lint`: replace `golangci-lint run ./...` → `go vet ./...`. In `check`: replace `golangci-lint run ./... && go vet ./...` → `go vet ./...`. Retry both. |
| `command not found: uvicorn` | `dev` (python) | Replace → `python -m src --reload` or skip with comment, retry |
| `command not found: ruff` | `lint`/`fmt`/`check` (python) | In `lint`: replace `ruff check .` → `python -m flake8`. In `check`: replace `ruff check .` → `python -m flake8` (keep `&& python -m py_compile src/`). In `fmt`: skip with comment. Retry. |

For each fix:
1. Edit the justfile to apply the correction.
2. Re-run the failed command (actual execution, same method as Phase 2).
3. If it still fails after 2 attempts, leave the recipe as-is and report the failure in the output.

#### 4d. Report verification results

After all recipes have been verified (or corrected):

```
Verification results:
  ✓ project-type    → "mixed" (executed)
  ✓ compile         → go vet ./... + npx tsc --noEmit (executed)
  ✓ build           → go build ./... + npm run build (executed)
  ✓ test            → go test ./... + npm test (dry-run only)
  ✗ run             → FIXED: npm start → npm run preview (executed, self-corrected)
  ✓ dev             → go run cmd/server/main.go + npm run dev (executed, 10s timeout)
  ✓ install         → go mod download + npm install (executed)
  ✗ lint            → golangci-lint not found, replaced with go vet (executed, self-corrected)

2 issues auto-corrected. Edit justfile to customize further.
```

### Step 5: Output Confirmation

```
Created justfile with standard forge targets (Go project)

Targets:
  just project-type               → echo "backend"
  just compile                    → go vet ./...
  just test                       → go test ./...
  just test-e2e --feature <slug>  → feature tests in tests/e2e/features/<slug>/
  ... (all 15 standard targets listed with resolved commands)

Edit justfile to customize commands for your project.
task all-completed will now use `just test` automatically.
```

## Notes

- **just >= 1.50.0**: `[arg("feature", long)]` generates `--feature <value>` named option syntax; callers (CI, `task all-completed`) must pass the slug: `just test-e2e --feature <slug>`
- Makefile migration: preserve original command logic, adjust only format
- **e2e tests are profile-aware**: The `test-e2e`, `e2e-setup`, and `e2e-verify` recipes are generated from the active test profile's `justfile-recipes` file (e.g., `plugins/forge/profiles/web-playwright/justfile-recipes`). Each profile defines its own execution commands. For `web-playwright`, this still uses `npx playwright test` and requires Node.js.
- **Targets invoked by forge skills**: `project-type`, `compile`, `build`, `test`, `test-e2e`, `install`, `e2e-setup`, `e2e-verify`. The remaining targets (`run`, `dev`, `lint`, `fmt`, `check`, `clean`, `ci`) are for manual use and are not called by any skill.
- **Idempotency**: `e2e-setup` and `install` are designed to be idempotent (safe to run multiple times). Other recipes (`build`, `compile`, `test`) are not — they always re-execute.
- **Mixed project scope**: forge skills resolve scope from `task claim` output or `process/state.json` and pass it to `just <verb>` when `just project-type` returns `mixed`. Pass `just compile frontend` or `just compile backend` manually to target a single side outside of a task context.

<EXTREMELY-IMPORTANT>
- MANUAL-ONLY. Do NOT auto-invoke this skill from other skills or agents. Only invoke when user explicitly runs `/init-justfile`.
- If an existing justfile lacks forge boundary markers and `--force` is not set, you MUST prompt the user before overwriting. Never silently destroy user customizations.
- Only the section between `# --- forge standard recipes ---` / `# --- end forge standard recipes ---` markers may be replaced. Recipes outside markers must be preserved verbatim.
- After writing, you MUST run the verification steps (dry-run + actual execution) and report all results.
</EXTREMELY-IMPORTANT>
