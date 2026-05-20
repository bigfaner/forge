---
name: init-justfile
description: Scaffold a Justfile with standard forge targets for the current project.
allowed-tools: Bash Read Write Edit
disable-model-invocation: true
argument-hint: '[--type frontend|backend|mixed] [--force]'
---

# /init-justfile

MANUAL-ONLY. Do NOT auto-invoke — only when user explicitly asks to `/init-justfile`.

Generate a Justfile with standard forge targets as an abstraction layer for test/build commands.

## Prerequisites

**Install just (>= 1.50.0)**

| Platform          | Command                                  |
| ----------------- | ---------------------------------------- |
| macOS / Linux     | `brew install just`                      |
| Windows (Scoop)   | `scoop install just`                     |
| Windows (winget)  | `winget install --id Casey.Just --exact` |
| Cargo (universal) | `cargo install just`                     |

```bash
just --version  # requires >= 1.50.0 (supports [arg] named option syntax)
```

If version < 1.50.0: `cargo install just`

## Parameters

| Parameter | Values                         | Default       | Description                                      |
| --------- | ------------------------------ | ------------- | ------------------------------------------------ |
| `--type`  | `frontend`, `backend`, `mixed` | (auto-detect) | Override project type                            |
| `--force` | (flag)                         | false         | Overwrite existing justfile without confirmation |

## Standard Target Contract

| Target         | Required | Purpose                                                                                           |
| -------------- | -------- | ------------------------------------------------------------------------------------------------- |
| `compile`      | No       | Type-check and transpile for fast feedback                                                        |
| `build`        | No       | Full compile and package                                                                          |
| `run`          | No       | Start the service                                                                                 |
| `dev`          | No       | Hot-reload development mode                                                                       |
| `test`         | Yes      | Unit + integration tests                                                                          |
| `test-e2e`     | No       | E2E tests; `--feature <slug>` for single feature                                                  |
| `lint`         | No       | Static analysis                                                                                   |
| `fmt`          | No       | Auto-format code                                                                                  |
| `check`        | No       | lint + compile (CI gate)                                                                          |
| `clean`        | No       | Remove build artifacts                                                                            |
| `install`      | No       | Install dependencies (idempotent)                                                                 |
| `ci`           | No       | Full CI pipeline                                                                                  |
| `e2e-setup`    | No       | Install e2e dependencies (idempotent)                                                             |
| `e2e-compile`  | No       | Compile-check e2e tests without running                                                           |
| `e2e-verify`   | No       | Check for unresolved `// VERIFY:` markers                                                         |

## Workflow

```
0. Load Convention → 1. Detect project type + entry points → 2. Check existing justfile → 3. Generate e2e recipes + assemble and write → 4. Verify and self-correct → 5. Output confirmation
```

### Step 0: Load Convention

Load test framework knowledge from Convention files. This provides the information needed to generate e2e recipes in Step 3.

1. List files in `docs/conventions/` directory.
2. For each file with `domains` frontmatter containing `testing`, read the file.
3. Extract the **Framework** section from the loaded Convention files.
4. From the Framework section, note:
   - Framework name (e.g., "Go testing package + testify/assert", "Vitest", "Ginkgo v2 + Gomega")
   - File pattern (e.g., `*_test.go`, `*.test.ts`)
   - Test runner (e.g., `go test`, `vitest run`, `ginkgo`)
   - Build tag / marker (e.g., `//go:build e2e`)
   - Result format output flags (e.g., `-json -v`, `--reporter=json`)
5. Also extract the **Tags** section for build-tag/marker syntax.
6. Also extract the **Result Format** section for execution command patterns.

<HARD-RULE>
Do NOT use framework-specific recipe templates. Generate e2e recipes from Convention content and LLM knowledge of the framework. The LLM constructs recipes based on the Convention description, not from hardcoded templates.
</HARD-RULE>

**If no Convention files found** (cold start):
- Proceed to Step 1 for file signal detection.
- LLM will generate e2e recipes from common patterns for the detected language/framework.
- Output hint: "No test Convention files found in docs/conventions/. Recipes will use LLM defaults. Run `/forge:test-guide` to create a Convention file."

### Step 1: Detect Project Type and Entry Points

```bash
just --version 2>/dev/null | awk '{print $2}' | awk -F. '$1 > 1 || ($1 == 1 && $2 >= 50)' | grep -q .
# If exit code != 0: output "Error: just >= 1.50.0 required — run `cargo install just`" and abort.
```

#### Parameter override

If `--type` is provided, skip project type detection (1a). Entry point detection (1b/1c) still runs for the detected language.

**Neither set**: Full auto-detection.

#### 1a. Project type detection

```bash
ls package.json go.mod Cargo.toml pyproject.toml 2>/dev/null
```

**Detection signal mapping:**

| Marker File      | Signal   |
| ---------------- | -------- |
| `package.json`   | frontend |
| `go.mod`         | backend  |
| `Cargo.toml`     | backend  |
| `pyproject.toml` | backend  |

**Classification algorithm:**

1. Check for each marker file's existence in the project root.
2. Count frontend vs backend signals and classify:
   - Exactly one frontend signal AND exactly one backend signal -> **`mixed`**
   - Exactly one frontend signal, no backend signals -> **`frontend`**
   - Exactly one backend signal, no frontend signals -> **`backend`**
   - Neither set has signals -> **Error**: "Error: no known project markers detected (expected one of: package.json, go.mod, Cargo.toml, pyproject.toml)" -- abort, do NOT generate a justfile.
   - Multiple backend signals (e.g. `go.mod` + `Cargo.toml`) -> **Error**: "Error: multiple backend markers detected -- not supported" -- abort.

**For `mixed` projects, also detect root paths:**

```bash
find . -name package.json -not -path '*/node_modules/*' -maxdepth 3 | head -1 | xargs dirname
find . \( -name go.mod -o -name Cargo.toml -o -name pyproject.toml \) -maxdepth 3 | head -1 | xargs dirname
```

Record these as `FRONTEND_DIR` and `BACKEND_DIR` (e.g. `./frontend`, `./backend`). Use `.` if the marker is in the project root.

#### 1b. Backend entry point detection

For **backend** and **mixed** projects, detect the entry point:

| Language | Marker           | Entry point detection (`BACKEND_ENTRY`)                                           |
| -------- | ---------------- | --------------------------------------------------------------------------------- |
| Go       | `go.mod`         | `ls cmd/*/main.go` -> `cmd/<name>/main.go`; else `ls main.go` -> `.`              |
| Rust     | `Cargo.toml`     | `grep '\[\[bin\]\]' Cargo.toml` -> `--bin <name>`; else empty                      |
| Python   | `pyproject.toml` | `ls src/__init__.py` -> `-m src`; `ls main.py` -> `main.py`; `ls app.py` -> `app.py` |

Record `BACKEND_ENTRY` from the detected entry point.

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
  - Check if it already contains `e2e-compile`, `e2e-test`, and `e2e-setup` recipes:
    ```bash
    just --list 2>/dev/null | grep -E 'e2e-compile|e2e-test|e2e-setup'
    ```
  - **If all three e2e recipes exist**: Output "justfile already contains e2e recipes (e2e-compile, e2e-test, e2e-setup). Skipping e2e recipe generation." Proceed to Step 4 for verification only.
  - **If some e2e recipes are missing**: Proceed to Step 3 to append only the missing recipes.
  - Check for boundary markers (`# --- forge standard recipes ---` / `# --- end forge standard recipes ---`).
  - **If boundary markers exist**: proceed to Step 3 (boundary marker merge).
  - **If boundary markers do NOT exist** (user's justfile has no forge markers):
    - If `--force` flag was provided: skip confirmation, proceed to Step 3.
    - If `--force` flag was NOT provided: prompt the user: "A justfile already exists without forge markers. Overwrite? (y/n)". If user declines, abort without modifying the file.
- If no justfile exists: proceed to Step 3 (create new file).

### Step 3: Generate e2e Recipes and Assemble Justfile

This step generates e2e-compile, e2e-test, and e2e-setup recipes from Convention knowledge and LLM understanding of the detected framework.

#### 3a. Generate e2e recipes from Convention

Using the Convention Framework, Tags, and Result Format sections loaded in Step 0:

**e2e-compile recipe**: Generate a recipe that compiles/checks e2e tests without running them. The recipe body is derived from the Convention's test runner and build tag info:

- Convention provides test runner (e.g., `go test`) and build tag (e.g., `//go:build e2e`) -> construct compile-check command (e.g., `go test -c ./tests/e2e/... -tags=e2e -o /dev/null` or `go vet -tags=e2e ./tests/e2e/...`)
- Convention provides file pattern and test runner (e.g., `vitest run`) -> construct type-check command (e.g., `npx tsc --noEmit` or `vitest run --passWithNoTests --run false`)

**e2e-test recipe**: Generate a recipe that runs e2e tests. The recipe body is derived from the Convention's Result Format execution command:

- Convention provides execution command (e.g., `go test ./tests/e2e/... -v -tags=e2e -json`) -> use as recipe body
- Convention provides test runner and output flags (e.g., `vitest run --reporter=json`) -> construct test command

**e2e-setup recipe**: Generate a recipe that installs e2e test dependencies:

- Convention indicates Go -> `go mod download`
- Convention indicates Node/Vitest -> `npm install` or `npx playwright install`
- Convention indicates Python -> `pip install -e ".[test]"` or `pip install pytest`
- Multiple frameworks in mixed project -> combine setup steps

**If no Convention was loaded** (cold start): LLM generates e2e recipes based on file signals detected in Step 1:

- `go.mod` detected -> generate Go e2e recipes using `go test` with `-tags=e2e`
- `package.json` detected -> generate Node e2e recipes using `npx vitest run` or appropriate test runner
- `Cargo.toml` detected -> generate Rust e2e recipes using `cargo test`
- `pyproject.toml` detected -> generate Python e2e recipes using `pytest`

The LLM uses conservative strategies for cold start recipes (most common patterns per language).

#### 3b. Generate non-e2e recipes

Generate standard project recipes (compile, build, test, run, dev, lint, fmt, check, clean, install, ci) based on detected language and project type from Step 1. The LLM constructs these from its knowledge of common tooling per language:

| Language | compile          | test            | lint                   | fmt               |
| -------- | ---------------- | --------------- | ---------------------- | ----------------- |
| Go       | `go vet ./...`   | `go test ./...` | `golangci-lint run ./...` | `gofmt -w .`      |
| Rust     | `cargo check`    | `cargo test`    | `cargo clippy -- -D warnings` | `cargo fmt` |
| Python   | `python -m py_compile src/` | `pytest` | `ruff check .` | `ruff format .` |
| Node     | `npx tsc --noEmit` | `npm test`   | `npx eslint .`         | `npx prettier --write .` |

For **mixed** projects, generate recipes with scope parameter:

```just
compile scope="":
    #!/usr/bin/env bash
    if [ -n "{{ scope }}" ]; then
        case "{{ scope }}" in
            backend)  {{ BACKEND_COMPILE }} ;;
            frontend) {{ FRONTEND_COMPILE }} ;;
            *)        {{ BACKEND_COMPILE }} && {{ FRONTEND_COMPILE }} ;;
        esac
    else
        {{ BACKEND_COMPILE }} && {{ FRONTEND_COMPILE }}
    fi
```

#### 3c. Placeholder substitution

For **single-language** projects:

| Placeholder    | Scope       | Replaced with                      |
| -------------- | ----------- | ---------------------------------- |
| `ENTRY_POINT`  | Go/Rust     | `BACKEND_ENTRY` from Step 1b       |
| `DEV_COMMAND`  | Python only | Dev server command from Step 1b    |
| `ENTRY_SCRIPT` | Node only   | `FRONTEND_RUN_SCRIPT` from Step 1c |

For **mixed** projects, replace `FRONTEND_DIR`, `BACKEND_DIR`, and all `BACKEND_*`/`FRONTEND_*` placeholders.

#### 3d. Boundary marker merge

When markers exist (`# --- forge standard recipes ---` / `# --- end forge standard recipes ---`), replace everything between them (inclusive) with the new recipes, preserving user recipes outside. Otherwise write the full template as a new file.

If justfile exists and is missing only some e2e recipes: append only the missing recipes without touching existing content.

### Step 4: Verify and Self-Correct

Two-phase verification: `--dry-run` catches syntax/structure errors, actual execution catches runtime errors.

#### 4a. Phase 1 -- Dry-run (syntax check)

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
just --dry-run e2e-compile
```

For mixed projects, also verify `--dry-run compile` and `--dry-run run` with `backend`/`frontend` scope arguments.

Fix any syntax failures before proceeding to Phase 2.

#### 4b. Phase 2 -- Actual execution (runtime check)

Execute each recipe for real to catch runtime errors. Recipes are classified by execution safety:

| Category                                              | Recipes                    | Method                                                                                    |
| ----------------------------------------------------- | -------------------------- | ----------------------------------------------------------------------------------------- |
| **Safe** (fast, no side effects)                      | `compile`, `lint`, `check` | Execute directly                                                                          |
| **Destructive** (modifies files or creates artifacts) | `build`, `fmt`, `clean`    | Execute directly (artifacts can be cleaned; fmt changes are welcome)                      |
| **Idempotent** (installs dependencies)                | `install`, `e2e-setup`     | Execute directly                                                                          |
| **Long-running** (starts servers)                     | `run`, `dev`               | Execute with timeout (10s), kill after timeout -- success = process still alive at timeout |
| **Expensive** (runs full test suite)                  | `test`, `test-e2e`         | Skip actual execution; verified by `--dry-run` only                                       |

For long-running recipes (`run`, `dev`): execute via `timeout 10 just <recipe> 2>&1 || true`. A crash before timeout ("missing script", "can't load package") is a runtime failure.

#### 4c. Self-correction rules

When a recipe fails in Phase 2, analyze the error and apply corrections:

| Error Pattern                         | Recipe                        | Fix                                                                                                                                                                                        |
| ------------------------------------- | ----------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `npm error Missing script: "start"`   | `run` (node/mixed)            | Replace `npm run start` -> `npm run preview` in justfile, retry                                                                                                                             |
| `npm error Missing script: "preview"` | `run` (node/mixed)            | Replace -> `npm run dev` in justfile, retry                                                                                                                                                 |
| `npm error Missing script: "dev"`     | `dev` (node/mixed)            | Replace -> `npm run start` in justfile, retry                                                                                                                                               |
| `can't load package: no Go files`     | `run`/`dev`/`compile` (go)    | Scan for `cmd/*/main.go`, update entry point in justfile, retry                                                                                                                            |
| `CGO_ENABLED=1` available             | `test` (go)                   | Add `-race` flag to `go test` recipe for race detection, retry                                                                                                                             |
| `command not found: golangci-lint`    | `lint`/`check` (go)           | In `lint`: replace `golangci-lint run ./...` -> `go vet ./...`. In `check`: replace `golangci-lint run ./... && go vet ./...` -> `go vet ./...`. Retry both.                                 |
| `command not found: uvicorn`          | `dev` (python)                | Replace -> `python -m src --reload` or skip with comment, retry                                                                                                                             |
| `command not found: ruff`             | `lint`/`fmt`/`check` (python) | In `lint`: replace `ruff check .` -> `python -m flake8`. In `check`: replace `ruff check .` -> `python -m flake8` (keep `&& python -m py_compile src/`). In `fmt`: skip with comment. Retry. |

For each fix:

1. Edit the justfile to apply the correction.
2. Re-run the failed command (actual execution, same method as Phase 2).
3. If it still fails after 2 attempts, leave the recipe as-is and report the failure in the output.

#### 4d. Report verification results

After all recipes have been verified (or corrected):

```
Verification results:
  [ok] compile         -> go vet ./... (executed)
  [ok] build           -> go build ./... (executed)
  [ok] test            -> go test ./... (dry-run only)
  [ok] e2e-compile     -> go test -c ./tests/e2e/... -tags=e2e -o /dev/null (dry-run only)
  [ok] e2e-test        -> go test ./tests/e2e/... -v -tags=e2e -json (dry-run only)
  [ok] e2e-setup       -> go mod download (executed)
  [fix] run            -> FIXED: updated entry point (executed, self-corrected)
  [ok] dev             -> go run cmd/server/main.go (executed, 10s timeout)
  [ok] install         -> go mod download (executed)
  [fix] lint           -> golangci-lint not found, replaced with go vet (executed, self-corrected)

2 issues auto-corrected. Edit justfile to customize further.
```

### Step 5: Output Confirmation

```
Created justfile with standard forge targets (Go project)

Targets:
  just compile                    -> go vet ./...
  just test                       -> go test ./...
  just e2e-compile                -> go test -c ./tests/e2e/... -tags=e2e -o /dev/null
  just e2e-test                   -> go test ./tests/e2e/... -v -tags=e2e -json
  just test-e2e --feature <slug>  -> feature tests in tests/e2e/features/<slug>/
  just e2e-setup                  -> go mod download
  ... (all standard targets listed with resolved commands)

Convention: docs/conventions/testing-go.md (Go testing package + testify/assert)
Edit justfile to customize commands for your project.
forge quality-gate will now use `just test` automatically.
```

If no Convention was used:

```
Created justfile with standard forge targets (Go project)

Targets:
  ... (all standard targets listed with resolved commands)

No Convention file found. Recipes generated from LLM defaults.
Run `/forge:test-guide` to create a Convention file for consistent future generation.
```

## Notes

- **just >= 1.50.0**: `[arg("feature", long)]` generates `--feature <value>` named option syntax; callers (CI, `forge quality-gate`) must pass the slug: `just test-e2e --feature <slug>`
- Makefile migration: preserve original command logic, adjust only format
- **Convention-driven e2e recipes**: The e2e-compile, e2e-test, and e2e-setup recipes are generated from Convention Framework/Tags/Result Format sections. No hardcoded framework templates. The LLM constructs recipes based on Convention content.
- **Cold start**: When no Convention files exist, the LLM generates recipes from common patterns for the detected language. These recipes use conservative defaults and may need manual adjustment.
- **Targets invoked by forge skills**: `compile`, `build`, `test`, `test-e2e`, `install`, `e2e-setup`, `e2e-compile`, `e2e-verify`. The remaining targets (`run`, `dev`, `lint`, `fmt`, `check`, `clean`, `ci`) are for manual use and are not called by any skill.
- **Idempotency**: `e2e-setup` and `install` are designed to be idempotent (safe to run multiple times). Other recipes (`build`, `compile`, `test`) are not -- they always re-execute.
- **Mixed project scope**: forge skills resolve scope from `forge task claim` output or `process/state.json` and pass it to `just <verb>`. Pass `just compile frontend` or `just compile backend` manually to target a single side outside of a task context.

<EXTREMELY-IMPORTANT>
- MANUAL-ONLY. Do NOT auto-invoke this skill from other skills or agents. Only invoke when user explicitly runs `/init-justfile`.
- If an existing justfile lacks forge boundary markers and `--force` is not set, you MUST prompt the user before overwriting. Never silently destroy user customizations.
- Only the section between `# --- forge standard recipes ---` / `# --- end forge standard recipes ---` markers may be replaced. Recipes outside markers must be preserved verbatim.
- After writing, you MUST run the verification steps (dry-run + actual execution) and report all results.
- Do NOT use framework-specific recipe templates. Generate e2e recipes from Convention content and LLM knowledge only.
</EXTREMELY-IMPORTANT>
