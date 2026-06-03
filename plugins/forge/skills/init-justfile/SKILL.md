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

`<prefix>` is `<key>-` for named surfaces (e.g., `app-test`) or empty for scalar surfaces (e.g., `test`).

| Target                  | Required | Purpose                                                              |
| ----------------------- | -------- | -------------------------------------------------------------------- |
| `<prefix>test-setup`    | No*      | Mobile only: prepare emulator and test environment                   |
| `<prefix>dev`           | No*      | Start dev server for surface (web/api/mobile only)                   |
| `<prefix>probe`         | No*      | Health check for surface (web/api/mobile only)                       |
| `<prefix>test`          | Yes      | Surface-level tests (functional for cli/tui/api, e2e for web/mobile) |
| `<prefix>teardown`      | Yes      | Stop services and cleanup                                            |
| `<prefix>`              | No*      | Aggregate: test-setup->dev->probe->test->teardown (mobile) / dev->probe->test->teardown (web/api) |

*CLI/TUI surfaces do NOT generate dev, probe, or aggregate recipes. Mobile also generates test-setup.

Each surface recipe MUST support `[linux]` and `[windows]` dual-platform variants.

### Recipe Parameter Signatures

| Recipe | Signature | Description |
|--------|-----------|-------------|
| `unit-test` | `just unit-test` (no parameters) | Language-level unit tests, no filtering needed |
| `<prefix>test` | `just <prefix>test [journey]` | Surface-level advanced tests; `<prefix>` is `<key>-` for named surfaces or empty for scalar. Optional journey parameter filters to a specific journey, omit to run all |
| `<prefix>probe` | `just <prefix>probe` (no parameters) | Surface health check |

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

Load test strategy knowledge from surface-first Convention files. This provides the information needed to generate `unit-test` and surface-level test recipes in Step 3.

**Surface-first loading** (new structure):

1. List subdirectories under `docs/conventions/testing/`.
2. For each subdirectory matching a configured surface key (from Step 1s):
   - Read `docs/conventions/testing/<surface>/core.md`.
   - Extract the **Assertion Preferences** table for framework-specific test runner, assertion library, and mock mechanism.
   - Extract the **File Location** section for test file patterns and directory rules.
   - Extract the **Lifecycle** section for build tag / marker syntax (e.g., `//go:build <surface>-<type>`).
3. From the collected Convention data, note per-surface:
   - Framework name (e.g., "Go testing package + testify/assert", "Vitest")
   - File pattern (e.g., `*_test.go`, `*.test.ts`)
   - Test runner (e.g., `go test`, `vitest run`)
   - Build tag / marker (e.g., `//go:build cli-functional`)
   - Result format output flags (e.g., `-json -v`, `--reporter=json`)

**Fallback — legacy structure**:

If `docs/conventions/testing/` contains flat files with `domains` frontmatter (e.g., `go.md`, `vitest.md`) instead of surface subdirectories:
1. Read each file with `domains` frontmatter containing `testing`.
2. Extract the **Framework**, **Tags**, and **Result Format** sections.
3. Output migration hint: "Legacy Convention structure detected (flat files). Run `/forge:test-guide` to migrate to surface-first structure (`testing/{surface}/`)."
4. Use extracted data as Convention knowledge for this run.

<HARD-RULE>
Use the language-specific template from `templates/<lang>.just` as the **starting point** for language-level recipe generation. Then customize recipes using Convention knowledge (from Step 0) and LLM understanding of the detected language/framework. The generation process is:

1. **Load template**: Read `templates/<lang>.just` (where `<lang>` is one of: `go`, `node`, `python`, `rust`, `mixed`). Use `templates/generic.just` as fallback when no language-specific template matches.
2. **Apply Convention overrides**: Replace recipe bodies with Convention-sourced commands where available (e.g., test runner, build tags, output flags from the Assertion Preferences and Lifecycle sections).
3. **LLM customization**: Adjust remaining recipes based on detected project structure (entry points, dependency managers, tool presence).

Templates define the **recipe structure** (names, groups, boundary markers, shared recipes like `probe` and `test-setup`). The LLM customizes **recipe bodies** (actual commands) based on Convention content and project signals.
</HARD-RULE>

**If no Convention files found** (cold start):
- Proceed to Step 1 for file signal detection.
- LLM will generate recipes from common patterns for the detected language/framework.
- Output hint: "No test Convention files found in docs/conventions/testing/. Recipes will use LLM defaults. Run `/forge:test-guide` to create surface-first Convention files."

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

Detect surfaces configured in the project:

```bash
forge surfaces 2>/dev/null
```

**Parsing rule**: Use the unified `forge surfaces` text output parsing rule (see Forge Guide → Surface Output Parsing).

**Outcome A — Surfaces configured** (text output is non-empty):

Parse each line using the rule above. Each line produces one surface entry with `key` and `type`:

| Text output | key | type | Form |
|-------------|-----|------|------|
| `tui` | (empty) | `tui` | scalar |
| `myapp=tui` | `myapp` | `tui` | named |
| `backend=api\nfrontend=web` | `backend` / `frontend` | `api` / `web` | named (multi) |

For each surface entry:
1. **Load surface rule file**: Read `rules/surfaces/<type>.md` relative to this SKILL.md's directory. If the file does not exist for the given type, emit warning: `Warning: no rule file for surface type "<type>" — skipping surface "<key>"` and skip this surface.

Collect the validated surfaces as `SURFACES_LIST`.

**Outcome B — No surfaces configured** (output is empty or command fails):

No surfaces detected. Skip surface recipe generation entirely. The justfile will contain only language-level targets — identical to the behavior before surface-aware support. This ensures **zero regression** for projects without surface configuration.

### Step 2: Check Existing Justfile

```bash
ls justfile Justfile 2>/dev/null
```

- If `justfile` or `Justfile` already exists:
  - Check if it already contains `unit-test`, `<prefix>test`, and `<prefix>teardown` recipes:
    ```bash
    just --list 2>/dev/null | grep -E 'unit-test|<prefix>test|<prefix>teardown'
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

This step generates language-level recipes (from language templates via Step 0 HARD-RULE + Convention knowledge) and surface-level recipes (from surface rule files).

#### 3a. Generate language-level recipes

Generate `unit-test`, `compile`, `build`, `lint`, `fmt`, `check`, `clean`, `install`, `ci` recipes using the three-layer approach defined in Step 0's HARD-RULE:

1. **Load language template**: Read the corresponding template file from `templates/<lang>.just`:
   - `go.just` for Go projects
   - `node.just` for Node/TypeScript projects
   - `python.just` for Python projects
   - `rust.just` for Rust projects
   - `mixed.just` for mixed frontend+backend projects
   - `generic.just` as fallback when no language-specific template matches

2. **Apply Convention overrides**: For each recipe body, check if Convention (loaded in Step 0) provides a more precise command:
   - `unit-test`: Use Convention's test runner + file pattern + output flags (e.g., `go test -json -v ./...` with `-tags` from Lifecycle section)
   - Other recipes: Override template defaults with Convention-sourced commands where available

3. **LLM customization**: Adjust remaining recipe bodies based on detected project structure:
   - Entry point detection from Step 1 (replace `ENTRY_POINT` / `ENTRY_SCRIPT` placeholders)
   - Dependency manager detection (npm vs pnpm vs yarn, pip vs poetry, etc.)
   - Tool presence checks (golangci-lint, prettier, ruff, etc.)

**unit-test recipe**: The template provides a language-appropriate default. Override with Convention's test runner + output flags if available. The `unit-test` recipe takes no parameters: `just unit-test`.

**Other language-level recipes** (compile, build, lint, fmt, check, clean, install, ci): Use template defaults as baseline, customize based on Convention and detected tooling.

**Template extras** (`run`, `dev`, `test`, `test-setup`, `probe`): Templates include additional convenience recipes beyond the Standard Target Contract. These are **project-level convenience targets** (not invoked by forge skills) and should be included in the generated justfile. Note: `test` and `probe` in templates are legacy pre-surface recipes — when surfaces are configured (Step 1s), these are superseded by surface-specific `<key>-test` and `<key>-probe` recipes generated in Step 3b.

#### 3b. Generate surface-level recipes

For each surface in `SURFACES_LIST` (collected in Step 1s):

1. Read the surface rule file loaded in Step 1s (`rules/surfaces/<type>.md`).
2. From the rule file, extract:
   - **Orchestration sequence**: which steps apply (dev/probe/test/teardown).
   - **Recipe contracts**: recipe names, signatures, exit codes.
   - **Journey filter strategy**: which journey tags belong to this surface.
3. Generate recipes for each step defined in the orchestration sequence. Recipe names use `<prefix><verb>` where `<prefix>` is `<key>-` for named surfaces or empty for scalar surfaces.

**Recipe naming** (determined by surface form from Step 1s parsing):
- **Scalar surface** (no key, only type): recipe names use the verb directly — `test`, `build`, `dev`, `teardown`. No prefix.
- **Named surface** (has key): recipe names use `<key>-<verb>` — e.g., `app-test`, `admin-panel-dev`. This applies to both single and multi-surface projects.

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
- Do NOT generate `<prefix>dev` or `<prefix>probe` recipes.
- Do NOT generate `<prefix>` aggregate recipe.
- Generate only `<prefix>test` and `<prefix>teardown`.

**Aggregate recipes** (web, api, mobile):
- Generate `<prefix>` aggregate recipe that calls sub-recipes in orchestration order.
- On any sub-step failure, run teardown and exit with the failure code:

```just
# scalar example (no key):
test-all:
    #!/usr/bin/env bash
    set -euo pipefail
    just dev && just probe && just test; rc=$?; just teardown; exit $rc

# named example (key=web):
web:
    #!/usr/bin/env bash
    set -euo pipefail
    just web-dev && just web-probe && just web-test; rc=$?; just web-teardown; exit $rc
```

**Recipe content generation**: The LLM fills in recipe bodies based on:
1. **Language template** (from Step 3a): provides the underlying command patterns (e.g., `go run`, `npm run dev`, `cargo run`) loaded from `templates/<lang>.just`.
2. **Convention knowledge** (from Step 0): provides framework-specific test runners and patterns to override template defaults.
3. **Surface rule file**: provides the orchestration sequence, exit code semantics, and recipe templates with `[linux]`/`[windows]` dual-platform structure.

The LLM synthesizes these three sources to produce concrete recipe bodies. The surface rule files provide TODO-stub templates as the structural skeleton; the LLM replaces stubs with actual commands derived from the language template and Convention knowledge. For example, a web surface in a Go project would generate:

```just
web-dev:
    go run ./cmd/server/main.go

web-probe:
    curl -sf http://localhost:8080/health

web-test:
    go test ./tests/... -v -tags=web-e2e -json

web-teardown:
    #!/usr/bin/env bash
    set -euo pipefail
    if [ -f tests/results/.pid-server ]; then
        kill "$(tr -d '\r' < tests/results/.pid-server)" 2>/dev/null || true
        rm -f tests/results/.pid-server
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

[group: <surface-key-or-type>]
<prefix>test-setup:   (mobile only)
<prefix>dev:          (if applicable)
<prefix>probe:        (if applicable)
<prefix>test:
<prefix>teardown:
<prefix>:              (aggregate, if applicable)

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

For surface recipes, verify each generated recipe (`<prefix>` is `<key>-` for named surfaces or empty for scalar):

```bash
just --dry-run <prefix>test-setup # (mobile only)
just --dry-run <prefix>dev       # (if applicable)
just --dry-run <prefix>probe     # (if applicable)
just --dry-run <prefix>test
just --dry-run <prefix>teardown
```

Fix any syntax failures before proceeding to Phase 2.

#### 4b. Phase 2 -- Actual execution (runtime check)

Execute each recipe for real to catch runtime errors. Recipes are classified by execution safety:

| Category                                              | Recipes                    | Method                                                                                    |
| ----------------------------------------------------- | -------------------------- | ----------------------------------------------------------------------------------------- |
| **Safe** (fast, no side effects)                      | `compile`, `lint`, `check` | Execute directly                                                                          |
| **Destructive** (modifies files or creates artifacts) | `build`, `fmt`, `clean`    | Execute directly (artifacts can be cleaned; fmt changes are welcome)                      |
| **Idempotent** (installs dependencies)                | `install`                  | Execute directly                                                                          |
| **Long-running** (starts servers)                     | `<prefix>dev`                | Execute with timeout (10s), kill after timeout -- success = process still alive at timeout |
| **Expensive** (runs full test suite)                  | `unit-test`, `<prefix>test`  | Skip actual execution; verified by `--dry-run` only                                       |

For long-running recipes (`<prefix>dev`): execute via `timeout 10 just <prefix>dev 2>&1 || true`. A crash before timeout ("missing script", "can't load package") is a runtime failure.

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
  [ok] web-test        -> go test ./tests/... -v -tags=web-e2e (dry-run only)
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
  mobile-app (mobile): mobile-test-setup, mobile-dev, mobile-probe, mobile-test, mobile-teardown, mobile

Language targets:
  just compile        -> go vet ./...
  just unit-test      -> go test ./...
  just lint           -> golangci-lint run ./...
  ... (all language targets listed)

Surface targets:
  just dev            -> go run ./cmd/server/main.go     (scalar, no prefix)
  just probe          -> curl -sf http://localhost:8080/health
  just test           -> go test ./tests/... -v -tags=web-e2e -json
  just teardown       -> (cleanup script)
  ... (repeat for each surface)
  -- OR named --
  just web-dev        -> go run ./cmd/server/main.go
  just web-probe      -> curl -sf http://localhost:8080/health
  just web-test       -> go test ./tests/... -v -tags=web-e2e -json
  just web-teardown   -> (cleanup script)
  just web            -> aggregate: dev->probe->test->teardown
  ... (repeat for each surface)

Convention: docs/conventions/testing/cli/core.md (Go testing package + testify/assert)
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

Convention: docs/conventions/testing/cli/core.md (Go testing package + testify/assert)
Edit justfile to customize commands for your project.
Run `/forge:init-justfile` again after configuring surfaces in .forge/config.yaml to add surface-aware recipes.
```

If no Convention was used:

```
No Convention file found. Recipes generated from LLM defaults.
Run `/forge:test-guide` to create surface-first Convention files for consistent future generation.
```

## Notes

- **just >= 1.50.0**: supports `[arg]` named option syntax and `[linux]`/`[windows]` platform attributes; surface recipes use dual-platform variants.
- **Zero regression**: Projects without surface configuration receive exactly the same justfile as before this feature. No new recipes, no changed behavior.
- **Two-layer model**: `unit-test` is language-level (fast, per-task submit gate); `<prefix>test` is surface-level (functional tests for cli/tui/api, e2e tests for web/mobile). Forge is surface-agnostic -- it calls `just <prefix>test` based on task surface-key. Test type terminology follows the [Surface Test Type Model](../test-guide/references/test-type-model.md).
- **Recipe naming**: Scalar surfaces (no key) produce prefix-less recipes (`test`, `dev`, `teardown`). Named surfaces produce `<key>-` prefixed recipes (e.g., `admin-panel-dev`, `payment-api-test`). Multi-surface projects always use named keys, so each surface gets its own prefix.
- **Targets invoked by forge skills**: `compile`, `unit-test`, `<prefix>test`, `<prefix>teardown`, `install`. The remaining targets are for manual use.
- **Cold start**: When no Convention files exist, the LLM generates recipes from common patterns for the detected language. These recipes use conservative defaults and may need manual adjustment.

<EXTREMELY-IMPORTANT>
- MANUAL-ONLY. Do NOT auto-invoke this skill from other skills or agents. Only invoke when user explicitly runs `/init-justfile`.
- If an existing justfile lacks forge boundary markers and `--force` is not set, you MUST prompt the user before overwriting. Never silently destroy user customizations.
- Only the section between `# --- forge standard recipes ---` / `# --- end forge standard recipes ---` markers may be replaced. Recipes outside markers must be preserved verbatim.
- After writing, you MUST run the verification steps (dry-run + actual execution) and report all results.
- Use language-specific templates from `templates/<lang>.just` as the starting point for recipe generation. Customize with Convention overrides and LLM knowledge. See Step 0 HARD-RULE for the three-layer generation process.
- Surface data source is `forge surfaces` (text mode). Scalar surfaces (no `=` in output line) produce prefix-less recipes (`test`, `dev`, `teardown`). Named surfaces produce `<key>-` prefixed recipes.
- CLI/TUI surfaces MUST NOT generate dev, probe, or aggregate recipes.
- `# user-customized` marked recipes MUST be preserved during re-generation.
</EXTREMELY-IMPORTANT>
