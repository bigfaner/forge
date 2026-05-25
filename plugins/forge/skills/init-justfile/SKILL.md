---
name: init-justfile
description: Scaffold a Justfile with standard forge targets for the current project, with surface-aware recipe generation.
allowed-tools: Bash Read Write Edit
disable-model-invocation: true
argument-hint: '[--type frontend|backend|mixed] [--force]'
---

# /init-justfile

MANUAL-ONLY. Do NOT auto-invoke — only when user explicitly asks to `/init-justfile`.

Generate a Justfile with standard forge targets as an abstraction layer for test/build commands. When surfaces are configured, generate surface-aware recipes (dev/probe/test/teardown) differentiated by surface type.

## Prerequisites

**Install just (>= 1.50.0)**

| Platform          | Command                                  |
| ----------------- | ---------------------------------------- |
| macOS / Linux     | `brew install just`                      |
| Windows (Scoop)   | `scoop install just`                     |
| Windows (winget)  | `winget install --id Casey.Just --exact` |
| Cargo (universal) | `cargo install just`                     |

```bash
just --version  # requires >= 1.50.0 (supports [arg] named option syntax + [linux]/[windows] attributes)
```

If version < 1.50.0: `cargo install just`

## Parameters

| Parameter | Values                         | Default       | Description                                      |
| --------- | ------------------------------ | ------------- | ------------------------------------------------ |
| `--type`  | `frontend`, `backend`, `mixed` | (auto-detect) | Override project type                            |
| `--force` | (flag)                         | false         | Overwrite existing justfile without confirmation |

## Standard Target Contract

### Language-Level Targets (always present)

| Target         | Required | Purpose                                                                                           |
| -------------- | -------- | ------------------------------------------------------------------------------------------------- |
| `compile`      | No       | Type-check and transpile for fast feedback                                                        |
| `build`        | No       | Full compile and package                                                                          |
| `unit-test`    | Yes      | Language-level unit tests (fast feedback, per-task submit gate)                                   |
| `lint`         | No       | Static analysis                                                                                   |
| `fmt`          | No       | Auto-format code                                                                                  |
| `check`        | No       | lint + compile (CI gate)                                                                          |
| `clean`        | No       | Remove build artifacts                                                                            |
| `install`      | No       | Install dependencies (idempotent)                                                                 |
| `ci`           | No       | Full CI pipeline                                                                                  |

### Surface-Level Targets (when surfaces configured)

| Target                | Required | Purpose                                                              |
| --------------------- | -------- | -------------------------------------------------------------------- |
| `<key>-dev`           | No*      | Start dev server for surface (web/api/mobile only)                   |
| `<key>-probe`         | No*      | Health check for surface (web/api/mobile only)                       |
| `<key>-test`          | Yes      | Surface-level advanced tests (e2e/integration)                       |
| `<key>-teardown`      | Yes      | Stop services and cleanup                                            |
| `<key>`               | No*      | Aggregate: dev->probe->test->teardown (web/api/mobile only)          |

*CLI/TUI surfaces do NOT generate dev, probe, or aggregate recipes.

Each surface recipe MUST support `[linux]` and `[windows]` dual-platform variants.

### Recipe Parameter Signatures

| Recipe | Signature | Description |
|--------|-----------|-------------|
| `unit-test` | `just unit-test` (no parameters) | Language-level unit tests, no filtering needed |
| `<key>-test` | `just <key>-test` (no parameters) | Surface-level advanced tests |
| `<key>-probe` | `just <key>-probe` (no parameters) | Surface health check |

## Process Flow

```
0. Load Convention
   -> 1. Detect project type + entry points
   -> 1s. Detect surfaces
   -> 2. Check existing justfile
   -> 3. Generate recipes (language + surface) + assemble
   -> 4. Verify and self-correct
   -> 5. Output confirmation
```

### Step 0: Load Convention

Load test framework knowledge from Convention files. This provides the information needed to generate `unit-test` and surface-level test recipes in Step 3.

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
6. Also extract the **Result Format** section for output flags and format type.

<HARD-RULE>
Do NOT use framework-specific recipe templates. Generate `unit-test` and surface-level test recipes from Convention content and LLM knowledge of the framework. The LLM constructs recipes based on the Convention description, not from hardcoded templates.
</HARD-RULE>

**If no Convention files found** (cold start):
- Proceed to Step 1 for file signal detection.
- LLM will generate recipes from common patterns for the detected language/framework.
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

Detect project type and entry points per `rules/project-detection.md`. This covers:
- Marker file scanning and classification (frontend/backend/mixed)
- Mixed project root path detection (`FRONTEND_DIR`, `BACKEND_DIR`)
- Backend entry point detection (`BACKEND_ENTRY`)
- Frontend run script detection (`FRONTEND_RUN_SCRIPT`)

### Step 1s: Detect Surfaces

<HARD-RULE>
Surface-key naming MUST match `[a-zA-Z0-9_-]+`. If a surface-key contains invalid characters, emit an error and abort — do NOT generate recipes for that surface.
</HARD-RULE>

Detect surfaces configured in the project:

```bash
forge surfaces --json 2>/dev/null
```

**Outcome A — Surfaces configured** (JSON output is non-empty, not `[]`, not `{"error":"..."}`):

Parse the JSON to extract surface entries. Each entry has `key` and `type` fields:

```json
{"surfaces": [{"key": "admin-panel", "type": "web"}, {"key": "payment-api", "type": "api"}]}
```

For each surface entry:
1. **Validate surface-key**: Check that `key` matches `[a-zA-Z0-9_-]+`. If not, emit error: `Error: invalid surface-key "<key>": must match [a-zA-Z0-9_-]+` and skip this surface.
2. **Load surface rule file**: Read `rules/surfaces/<type>.md` relative to this SKILL.md's directory. If the file does not exist for the given type, emit warning: `Warning: no rule file for surface type "<type>" — skipping surface "<key>"` and skip this surface.

Collect the validated surfaces as `SURFACES_LIST`.

**Outcome B — No surfaces configured** (output is `[]`, or `{"error":"..."}`, or command fails):

No surfaces detected. Skip surface recipe generation entirely. The justfile will contain only language-level targets — identical to the behavior before surface-aware support. This ensures **zero regression** for projects without surface configuration.

### Step 2: Check Existing Justfile

```bash
ls justfile Justfile 2>/dev/null
```

- If `justfile` or `Justfile` already exists:
  - Check if it already contains `unit-test`, `<key>-test`, and `<key>-teardown` recipes:
    ```bash
    just --list 2>/dev/null | grep -E 'unit-test|<key>-test|<key>-teardown'
    ```
  - **If all expected recipes exist**: Output "justfile already contains required recipes. Skipping recipe generation." Proceed to Step 4 for verification only.
  - **If some recipes are missing**: Proceed to Step 3 to append only the missing recipes.
  - Check for boundary markers (`# --- forge standard recipes ---` / `# --- end forge standard recipes ---`).
  - **If boundary markers exist**: proceed to Step 3 (boundary marker merge).
  - **If boundary markers do NOT exist** (user's justfile has no forge markers):
    - If `--force` flag was provided: skip confirmation, proceed to Step 3.
    - If `--force` flag was NOT provided: prompt the user: "A justfile already exists without forge markers. Overwrite? (y/n)". If user declines, abort without modifying the file.
- If no justfile exists: proceed to Step 3 (create new file).

<HARD-RULE>
If an existing justfile lacks forge boundary markers and `--force` is not set, you MUST prompt the user before overwriting. Never silently destroy user customizations.
</HARD-RULE>

### Step 3: Generate Recipes and Assemble Justfile

This step generates language-level recipes (from language templates + Convention knowledge) and surface-level recipes (from surface rule files).

#### 3a. Generate language-level recipes

Generate `unit-test`, `compile`, `build`, `lint`, `fmt`, `check`, `clean`, `install`, `ci` recipes from Convention knowledge and LLM understanding of the detected language/framework.

**unit-test recipe**: Generate a language-level unit test recipe:

| Language | `unit-test` recipe body |
| -------- | ----------------------- |
| Go       | `go test ./...`         |
| Rust     | `cargo test`            |
| Python   | `pytest`                |
| Node     | `npm test`              |

- If Convention provides a specific test runner and file pattern, construct the appropriate command.
- The `unit-test` recipe takes no parameters: `just unit-test`.

**Other language-level recipes** (compile, build, lint, fmt, check, clean, install, ci): Generate based on detected language per Step 1. The LLM constructs these from its knowledge of common tooling per language.

#### 3b. Generate surface-level recipes

For each surface in `SURFACES_LIST` (collected in Step 1s):

1. Read the surface rule file loaded in Step 1s (`rules/surfaces/<type>.md`).
2. From the rule file, extract:
   - **Orchestration sequence**: which steps apply (dev/probe/test/teardown).
   - **Recipe contracts**: recipe names, signatures, exit codes.
   - **Journey filter strategy**: which journey tags belong to this surface.
3. Generate recipes for each step defined in the orchestration sequence.

**Recipe naming**:
- **Single surface project** (only one surface in `SURFACES_LIST`): use `<type>-<verb>` naming (e.g., `web-dev`, `web-test`). The surface-key IS the type.
- **Mixed project** (multiple surfaces): use `<surface-key>-<verb>` naming (e.g., `admin-panel-dev`, `payment-api-test`). This distinguishes recipes between surfaces.

**Dual-platform variants**: Each recipe MUST have `[linux]` and `[windows]` attribute variants:

```just
# user-customized
web-dev:
    #!/usr/bin/env bash
    set -euo pipefail
    # Linux-specific dev command

# user-customized
web-dev:
    #!/usr/bin/env bash
    set -euo pipefail
    # Windows-specific dev command
```

<HARD-RULE>
Every surface recipe MUST include `# user-customized` comment above the `[linux]` variant. This marks the recipe as user-customizable. Recipes with `# user-customized` markers in an existing justfile MUST NOT be overwritten during re-generation.
</HARD-RULE>

**Surface types without dev/probe** (cli, tui):
- Do NOT generate `<key>-dev` or `<key>-probe` recipes.
- Do NOT generate `<key>` aggregate recipe.
- Generate only `<key>-test` and `<key>-teardown`.

**Aggregate recipes** (web, api, mobile):
- Generate `<key>` aggregate recipe that calls sub-recipes in orchestration order.
- On any sub-step failure, run teardown and exit with the failure code:

```just
web:
    #!/usr/bin/env bash
    set -euo pipefail
    just web-dev && just web-probe && just web-test; rc=$?; just web-teardown; exit $rc
```

**Recipe content generation**: The LLM fills in recipe bodies based on:
1. **Language template** (from Step 1): determines the underlying command (e.g., `go run`, `npm run dev`, `cargo run`).
2. **Convention knowledge** (from Step 0): provides framework-specific test runners and patterns.
3. **Surface rule file**: provides the orchestration sequence and exit code semantics.

The LLM synthesizes these three sources to produce concrete recipe bodies. For example, a web surface in a Go project would generate:

```just
web-dev:
    go run ./cmd/server/main.go

web-probe:
    curl -sf http://localhost:8080/health

web-test:
    go test ./tests/e2e/... -v -tags=e2e -json

web-teardown:
    #!/usr/bin/env bash
    set -euo pipefail
    if [ -f tests/e2e/results/.pid-server ]; then
        kill "$(tr -d '\r' < tests/e2e/results/.pid-server)" 2>/dev/null || true
        rm -f tests/e2e/results/.pid-server
    fi
```

#### 3c. Boundary marker merge

When markers exist (`# --- forge standard recipes ---` / `# --- end forge standard recipes ---`), replace everything between them (inclusive) with the new recipes, preserving user recipes outside. Otherwise write the full template as a new file.

<HARD-RULE>
Only the section between `# --- forge standard recipes ---` / `# --- end forge standard recipes ---` markers may be replaced. Recipes outside markers MUST be preserved verbatim.
</HARD-RULE>

If justfile exists and is missing only some recipes: append only the missing recipes without touching existing content.

**User-customized protection**: Before replacing or removing any recipe, check if it has a `# user-customized` comment marker. If it does, preserve the existing recipe body unchanged — do NOT overwrite it with the generated version.

#### 3d. Recipe organization

Organize the justfile in this order:

```
# --- forge standard recipes ---

[group: language]
compile:
build:
install:
clean:
ci:

[group: language-test]
unit-test:
lint:
fmt:
check:

[group: <surface-key>]
<key>-dev:          (if applicable)
<key>-probe:        (if applicable)
<key>-test:
<key>-teardown:
<key>:              (aggregate, if applicable)

... (repeat for each surface)

# --- end forge standard recipes ---
```

### Step 4: Verify and Self-Correct

Two-phase verification: `--dry-run` catches syntax/structure errors, actual execution catches runtime errors.

#### 4a. Phase 1 -- Dry-run (syntax check)

Run each recipe with `--dry-run` to verify recipe syntax, variable expansion, and command structure:

```bash
just --list
just --dry-run compile
just --dry-run build
just --dry-run unit-test
just --dry-run lint
just --dry-run fmt
just --dry-run check
just --dry-run install
```

For surface recipes, verify each generated recipe:

```bash
just --dry-run <key>-dev       # (if applicable)
just --dry-run <key>-probe     # (if applicable)
just --dry-run <key>-test
just --dry-run <key>-teardown
```

Fix any syntax failures before proceeding to Phase 2.

#### 4b. Phase 2 -- Actual execution (runtime check)

Execute each recipe for real to catch runtime errors. Recipes are classified by execution safety:

| Category                                              | Recipes                    | Method                                                                                    |
| ----------------------------------------------------- | -------------------------- | ----------------------------------------------------------------------------------------- |
| **Safe** (fast, no side effects)                      | `compile`, `lint`, `check` | Execute directly                                                                          |
| **Destructive** (modifies files or creates artifacts) | `build`, `fmt`, `clean`    | Execute directly (artifacts can be cleaned; fmt changes are welcome)                      |
| **Idempotent** (installs dependencies)                | `install`                  | Execute directly                                                                          |
| **Long-running** (starts servers)                     | `<key>-dev`                | Execute with timeout (10s), kill after timeout -- success = process still alive at timeout |
| **Expensive** (runs full test suite)                  | `unit-test`, `<key>-test`  | Skip actual execution; verified by `--dry-run` only                                       |

For long-running recipes (`<key>-dev`): execute via `timeout 10 just <key>-dev 2>&1 || true`. A crash before timeout ("missing script", "can't load package") is a runtime failure.

#### 4c. Self-correction rules

When a recipe fails in Phase 2, analyze the error and apply corrections per `rules/self-correction.md`:

1. Match the error against known error patterns (npm missing scripts, Go package issues, missing linters/formatters).
2. Edit the justfile to apply the correction.
3. Re-run the failed command (actual execution, same method as Phase 2).
4. If it still fails after 2 attempts, leave the recipe as-is and report the failure in the output.

#### 4d. Report verification results

After all recipes have been verified (or corrected):

```
Verification results:
  [ok] compile         -> go vet ./... (executed)
  [ok] unit-test       -> go test ./... (dry-run only)
  [ok] web-dev         -> go run ./cmd/server/main.go (executed, 10s timeout)
  [ok] web-probe       -> curl -sf http://localhost:8080/health (executed)
  [ok] web-test        -> go test ./tests/e2e/... -v -tags=e2e (dry-run only)
  [ok] web-teardown    -> bash cleanup script (executed)
  [fix] lint           -> golangci-lint not found, replaced with go vet (executed, self-corrected)

1 issue auto-corrected. Edit justfile to customize further.
```

### Step 5: Output Confirmation

**With surfaces configured:**

```
Created justfile with surface-aware forge targets (Go project)

Surfaces:
  admin-panel (web): web-dev, web-probe, web-test, web-teardown, web
  payment-api (api): api-dev, api-probe, api-test, api-teardown, api

Language targets:
  just compile        -> go vet ./...
  just unit-test      -> go test ./...
  just lint           -> golangci-lint run ./...
  ... (all language targets listed)

Surface targets:
  just web-dev        -> go run ./cmd/server/main.go
  just web-probe      -> curl -sf http://localhost:8080/health
  just web-test       -> go test ./tests/e2e/... -v -tags=e2e -json
  just web-teardown   -> (cleanup script)
  just web            -> aggregate: dev->probe->test->teardown
  ... (repeat for each surface)

Convention: docs/conventions/testing/go.md (Go testing package + testify/assert)
Edit justfile to customize commands for your project.
Recipes marked `# user-customized` will be preserved on re-generation.
forge quality-gate will now use `just unit-test` for per-task gates.
```

**Without surfaces (zero-regression path):**

```
Created justfile with standard forge targets (Go project)

Targets:
  just compile                    -> go vet ./...
  just unit-test                  -> go test ./...
  just lint                       -> golangci-lint run ./...
  ... (all standard targets listed with resolved commands)

Convention: docs/conventions/testing/go.md (Go testing package + testify/assert)
Edit justfile to customize commands for your project.
Run `/forge:init-justfile` again after configuring surfaces in .forge/config.yaml to add surface-aware recipes.
```

If no Convention was used:

```
No Convention file found. Recipes generated from LLM defaults.
Run `/forge:test-guide` to create a Convention file for consistent future generation.
```

## Notes

- **just >= 1.50.0**: supports `[arg]` named option syntax and `[linux]`/`[windows]` platform attributes; surface recipes use dual-platform variants.
- **Surface-aware recipes**: When `.forge/config.yaml` defines surfaces, the justfile gains per-surface dev/probe/test/teardown recipes with `[linux]`/`[windows]` dual-platform variants. Without surfaces, the justfile is identical to the pre-surface-aware version.
- **Zero regression**: Projects without surface configuration receive exactly the same justfile as before this feature. No new recipes, no changed behavior.
- **Convention-driven test recipes**: Recipes are generated from Convention content + surface rule files + LLM knowledge. No hardcoded framework templates.
- **Two-layer model**: `unit-test` is language-level (fast, per-task submit gate); `<key>-test` is surface-level (advanced tests). Forge is surface-agnostic -- it calls `just <key>-test` based on task surface-key.
- **User-customized protection**: Recipes marked with `# user-customized` are never overwritten during re-generation. This allows users to customize surface recipes without fear of losing changes.
- **Mixed project naming**: When multiple surfaces exist, recipes use the surface-key as prefix (e.g., `admin-panel-dev`, `payment-api-test`) to avoid collisions. Single-surface projects use the surface-type as prefix (e.g., `web-dev`).
- **Targets invoked by forge skills**: `compile`, `unit-test`, `<key>-test`, `<key>-teardown`, `install`. The remaining targets are for manual use.
- **Cold start**: When no Convention files exist, the LLM generates recipes from common patterns for the detected language. These recipes use conservative defaults and may need manual adjustment.

<EXTREMELY-IMPORTANT>
- MANUAL-ONLY. Do NOT auto-invoke this skill from other skills or agents. Only invoke when user explicitly runs `/init-justfile`.
- If an existing justfile lacks forge boundary markers and `--force` is not set, you MUST prompt the user before overwriting. Never silently destroy user customizations.
- Only the section between `# --- forge standard recipes ---` / `# --- end forge standard recipes ---` markers may be replaced. Recipes outside markers must be preserved verbatim.
- After writing, you MUST run the verification steps (dry-run + actual execution) and report all results.
- Do NOT use framework-specific recipe templates. Generate test recipes from Convention content, surface rule files, and LLM knowledge only.
- Surface-key MUST match `[a-zA-Z0-9_-]+`. Abort on invalid keys — never generate recipes for invalid surface names.
- CLI/TUI surfaces MUST NOT generate dev, probe, or aggregate recipes.
- `# user-customized` marked recipes MUST be preserved during re-generation.
</EXTREMELY-IMPORTANT>
