---
name: init-justfile
description: Scaffold a Justfile with standard forge targets for the current project, with surface-aware recipe generation.
allowed-tools: Bash Read Write Edit
disable-model-invocation: true
argument-hint: '[--force]'
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

| Parameter | Values  | Default | Description                                      |
| --------- | ------- | ------- | ------------------------------------------------ |
| `--force` | (flag)  | false   | Overwrite existing justfile without confirmation |

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

#### Test Directory Path in Recipes

The `<prefix>test` recipe must resolve test scripts from the correct directory, which depends on the project's surface count:

- **Single surface** (1 surface, scalar or named): `tests/<journey>/` -- no surface key directory layer.
- **Multi surface** (2+ surfaces): `tests/<key>/<journey>/` -- each surface's tests live under its key.

When generating the `<prefix>test` recipe body, the agent must include the correct base path. For multi-surface projects, the `<key>` directory segment uses the surface's **key** (e.g., `backend`, `frontend`), NOT its type (e.g., `api`, `web`).

Example: In a project with `backend=api` + `frontend=web` surfaces, the `backend-test` recipe runs tests from `tests/backend/<journey>/`, and `frontend-test` from `tests/frontend/<journey>/`.

### Surface Gate Targets (multi-surface projects)

When a project has multiple surfaces (2+ named surfaces), the per-task quality gate validates only the surface affected by the task. This requires surface-specific gate recipes that mirror the language-level targets but scoped to a single surface.

Gate recipes use the same `<key>-` prefix convention as lifecycle recipes:

| Target                | Required | Purpose                                                           |
| --------------------- | -------- | ----------------------------------------------------------------- |
| `<key>-compile`       | No       | Type-check and transpile for the `<key>` surface code only        |
| `<key>-fmt`           | No       | Auto-format code for the `<key>` surface only                     |
| `<key>-lint`          | No       | Static analysis for the `<key>` surface code only                 |
| `<key>-unit-test`     | No       | Unit tests for the `<key>` surface code only (per-task submit gate) |

These recipes are **only generated for multi-surface projects** (2+ named surfaces). Single-surface projects use the language-level targets directly — no gate recipes are needed.

**Resolution logic**: At per-task gate time, if the task has a `surface-key` (e.g., `backend`), the gate resolves `<key>-<recipe>` (e.g., `backend-compile`). If the prefixed recipe exists, it is used; otherwise, the gate falls back to the generic language-level recipe (e.g., `compile`). This fallback ensures backward compatibility with projects that have surfaces but no prefixed gate recipes yet.

**Key vs Type**: Gate recipe prefixes use the surface **key** (e.g., `backend`, `frontend`), NOT the surface type (e.g., `api`, `web`). This matches the lifecycle recipe naming convention.

**Scope isolation**: Each gate recipe MUST validate ONLY its own surface. For example, `backend-compile` must not trigger frontend compilation, and `frontend-lint` must not lint backend code.

### Recipe Parameter Signatures

| Recipe | Signature | Description |
|--------|-----------|-------------|
| `unit-test` | `just unit-test` (no parameters) | Language-level unit tests, no filtering needed |
| `<prefix>test` | `just <prefix>test [journey]` | Surface-level advanced tests; optional journey parameter filters to a specific journey, omit to run all |
| `<prefix>probe` | `just <prefix>probe` (no parameters) | Surface health check |

## Process Flow

```
0. Detect surfaces (prerequisite)
   -> 1. Detect languages + load Convention
   -> 2. Check existing justfile
   -> 3. Generate recipes (agent-driven: language + surface) + assemble
   -> 4. Verify and self-correct
   -> 5. Output confirmation
```

### Step 0: Detect Surfaces

Surfaces are a prerequisite for surface-aware recipe generation. Detect surfaces configured in the project:

```bash
forge surfaces 2>/dev/null
```

**Parsing rule**: Use the unified `forge surfaces` text output parsing rule (see Forge Guide -> Surface Output Parsing).

**Outcome A — Surfaces configured** (text output is non-empty):

Parse each line using the rule above. Each line produces one surface entry with `key` and `type`:

| Text output | key | type | Form |
|-------------|-----|------|------|
| `tui` | (empty) | `tui` | scalar |
| `myapp=tui` | `myapp` | `tui` | named |
| `backend=api\nfrontend=web` | `backend` / `frontend` | `api` / `web` | named (multi) |

For each surface entry:
1. **Load surface rule file**: Read `rules/surfaces/<type>.md` relative to this SKILL.md's directory. If the file does not exist for the given type, emit warning: `Warning: no rule file for surface type "<type>" -- skipping surface "<key>"` and skip this surface.

Collect the validated surfaces as `SURFACES_LIST`.

**Outcome B — No surfaces configured** (output is empty or command fails):

No surfaces detected. Output the following message and abort:

```
No surfaces configured in this project. Run `forge init` to configure surfaces before generating a justfile.

Surface-aware recipe generation requires surface definitions (e.g., "backend=api", "frontend=web").
After running `forge init`, re-run `/init-justfile` to generate surface-aware recipes.
```

This ensures users are guided to set up surfaces rather than silently receiving a justfile without surface recipes.

### Step 1: Detect Languages and Load Convention

This step detects the language for each surface and loads framework-specific knowledge from Convention files. The combined knowledge drives agent-driven recipe generation in Step 3.

**Execution order**: First scan the project root (1a), then for each surface in `SURFACES_LIST` (from Step 0) scan its working directory independently (1b), then load Convention files per surface (1c).

#### 1a. Root Language Detection

Scan the project root for marker files to determine the primary language:

| Marker File               | Language     | Key Signals                                      |
| ------------------------- | ------------ | ------------------------------------------------ |
| `go.mod`                  | Go           | Module path, Go version                          |
| `package.json`            | Node/TypeScript | Scripts section, dependencies                 |
| `Cargo.toml`              | Rust         | Edition, dependencies                            |
| `pyproject.toml`/`setup.py` | Python    | Build system, dependencies                       |
| `pom.xml`/`build.gradle`  | Java        | Group/artifact, plugins                          |
| `build.gradle.kts`        | Kotlin       | Kotlin plugins, dependencies                     |
| `Gemfile`                 | Ruby         | Gems, Ruby version                               |
| `*.csproj`/`*.sln`        | C#/.NET      | Target framework, dependencies                   |
| `go.work`                 | Go (workspace) | Multiple Go modules                           |

**Detection process**:
1. Check the project root for marker files.
2. If a marker file exists, read it to extract language version, dependency manager, and key tooling (e.g., `npm` vs `pnpm` vs `yarn` from `package.json`'s `packageManager` field).
3. Record the detected language and tooling as the project's primary language.

If no marker file matches any known pattern, the agent falls back to its own knowledge to infer the language from file extensions and project structure.

#### 1b. Per-Surface Language Detection

For each surface in `SURFACES_LIST` (from Step 0), scan the surface's working directory for marker files using the same table in 1a. In multi-surface projects, different surfaces may use different languages (e.g., Go backend + Node frontend). Each surface records its own detected language independently.

For single-surface scalar projects, the surface's working directory is the project root — this step reuses the detection from 1a.

#### 1c. Load Convention Knowledge

Load test strategy and framework knowledge from Convention files. Convention files provide the information needed to generate `unit-test` and surface-level test recipes in Step 3.

**Surface-first loading** (preferred structure):

1. List subdirectories under `docs/conventions/testing/`.
2. For each subdirectory matching a surface key in `SURFACES_LIST`:
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

**If no Convention files found** (cold start):
- Proceed using agent's built-in knowledge of common patterns for the detected language/framework.
- Output hint: "No test Convention files found in docs/conventions/testing/. Recipes will use agent defaults. Run `/forge:test-guide` to create surface-first Convention files."

#### 1d. Load Server Lifecycle Patterns

<HARD-RULE>
When generating surface-level recipes that involve server lifecycle (dev, probe, teardown, and test recipes with embedded server startup), the agent MUST reference `rules/server-lifecycle.md` for ready-to-use bash code snippets. Prefer reusing these snippets (with slot placeholder substitution) over generating lifecycle code from scratch. This ensures reliability for PID tracking, idempotent start, health check, and teardown patterns.
</HARD-RULE>

The server lifecycle rule provides:
- PID file management (atomic write, stale detection, cleanup)
- Idempotent start (three-layer check: tracked process alive, port occupancy, start if needed)
- Health check (HTTP probe with retry, TCP probe)
- Multi-service orchestration (dependency order, per-service PID isolation)
- Complete recipe snippets for dev, teardown, and test-with-server-lifecycle

### Step 2: Check Existing Justfile

```bash
ls justfile Justfile 2>/dev/null
```

- If no justfile exists: proceed to Step 3 (create new file).
- If `justfile` or `Justfile` already exists, evaluate in this priority order:
  1. **Check boundary markers** (`# --- forge standard recipes ---` / `# --- end forge standard recipes ---`):
     - **Markers exist**: Check recipe completeness:
       ```bash
       just --list 2>/dev/null | grep -E 'unit-test|<prefix>test|<prefix>teardown'
       ```
       - All expected recipes present: Output "justfile already contains required recipes. Skipping recipe generation." Proceed to Step 4 for verification only.
       - Some recipes missing: Proceed to Step 3 (boundary marker merge, append missing recipes only).
     - **Markers do NOT exist**:
       - If `--force` flag was provided: skip confirmation, proceed to Step 3.
       - If `--force` flag was NOT provided: prompt the user: "A justfile already exists without forge markers. Overwrite? (y/n)". If user declines, abort without modifying the file.

<HARD-RULE>
If an existing justfile lacks forge boundary markers and `--force` is not set, you MUST prompt the user before overwriting. Never silently destroy user customizations.
</HARD-RULE>

### Step 3: Generate Recipes and Assemble Justfile

This step uses agent-driven generation: the agent synthesizes language-level and surface-level recipes from three knowledge sources, with no language templates.

**Knowledge source priority** (highest to lowest):
1. **Surface rule files** — recipe contracts (names, parameters, exit code semantics, orchestration sequence)
2. **Convention files** — framework-specific knowledge (test runner, build tool, lint tool, file patterns)
3. **Agent knowledge** — mainstream language/framework command patterns (cold start fallback)

#### 3a. Generate language-level recipes

Generate `unit-test`, `compile`, `build`, `lint`, `fmt`, `check`, `clean`, `install`, `ci` recipes based on the detected language (from Step 1a) and Convention knowledge (from Step 1c).

The agent determines the correct commands for each recipe by analyzing:

1. **Detected language and tooling** (from Step 1a/1b marker file scanning):
   - Language-specific build/test/lint commands (e.g., `go test`, `npm test`, `cargo test`, `pytest`)
   - Dependency manager (e.g., `npm` vs `pnpm` vs `yarn`, `pip` vs `poetry`, `cargo` vs `go mod`)
   - Tool presence (e.g., golangci-lint, prettier, ruff, eslint)

2. **Convention overrides** (from Step 1c):
   - `unit-test`: Use Convention's test runner + file pattern + output flags (e.g., `go test -json -v ./...` with `-tags` from Lifecycle section)
   - Other recipes: Use Convention-sourced commands where available

3. **Agent knowledge** (fallback when Convention is absent):
   - Mainstream defaults for the detected language (e.g., Go: `go build ./...`, `go vet ./...`; Node: `npm run build`, `npm run lint`)

**unit-test recipe**: Generate a language-appropriate command. Override with Convention's test runner + output flags if available. The `unit-test` recipe takes no parameters: `just unit-test`.

**Other language-level recipes** (compile, build, lint, fmt, check, clean, install, ci): Generate based on detected tooling and Convention knowledge. Each recipe must have `[linux]` and `[windows]` dual-platform variants where commands differ between platforms.

**Error stub fallback**: If the agent cannot determine a meaningful command for a recipe (e.g., rare language without known tooling), generate an error stub: `echo "TODO: implement <recipe-name>" >&2; exit 1`. This matches the behavior of the previous generic template.

#### 3b. Generate surface-level recipes

For each surface in `SURFACES_LIST` (collected in Step 0):

1. Read the surface rule file loaded in Step 0 (`rules/surfaces/<type>.md`).
2. From the rule file, extract:
   - **Orchestration sequence**: which steps apply (dev/probe/test/teardown).
   - **Recipe contracts**: recipe names, signatures, exit codes.
   - **Journey filter strategy**: which journey tags belong to this surface.
3. Generate recipes for each step defined in the orchestration sequence. Recipe names use `<prefix><verb>` (see `<prefix>` definition in Standard Target Contract).

**Recipe naming** (determined by surface form from Step 0 parsing):
- **Scalar surface** (no key, only type): recipe names use the verb directly — `test`, `build`, `dev`, `teardown`. No prefix.
- **Named surface** (has key): recipe names use `<key>-<verb>` — e.g., `app-test`, `admin-panel-dev`. This applies to both single and multi-surface projects.

**Dual-platform variants**: Each recipe MUST have `[linux]` and `[windows]` attribute variants:

```just
# user-customized
web-dev [linux]:
    #!/usr/bin/env bash
    set -euo pipefail
    # Linux-specific dev command

# user-customized
web-dev [windows]:
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
- On any sub-step failure, run teardown and exit with the failure code.

Aggregate pattern (web/api): `just <prefix>dev && just <prefix>probe && just <prefix>test; rc=$?; just <prefix>teardown; exit $rc`

Note: mobile surface uses a different aggregate pattern that includes `test-setup` as the first step — see the mobile surface rule file for the correct sequence.

**Test directory path**: The `<prefix>test` recipe body must resolve test scripts from the correct path (see "Test Directory Path in Recipes" above).

**Recipe content generation**: The agent fills in recipe bodies by synthesizing:
1. **Detected language and tooling** (from Step 1a/1b): provides the underlying command patterns (e.g., `go run`, `npm run dev`, `cargo run`).
2. **Convention knowledge** (from Step 1c): provides framework-specific test runners and patterns.
3. **Surface rule file**: provides the orchestration sequence, exit code semantics, and recipe naming contract.
4. **Server lifecycle patterns** (from `rules/server-lifecycle.md`, Step 1d): provides ready-to-use bash snippets for dev, probe, teardown, and test-with-server-lifecycle recipes. The agent replaces slot placeholders with project-specific values.

Example (web surface, Go project): `web-dev` uses idempotent start from server-lifecycle.md with `go run` as `<START_CMD>`, `web-probe` uses HTTP probe pattern, `web-test` uses test-with-server-lifecycle pattern with `go test -tags=web-e2e ./tests/<key>/{{journey}}/...` as `<TEST_CMD>`, `web-teardown` uses PID cleanup pattern.

**Slot placeholder resolution** (for server-lifecycle.md snippets): When Convention files are absent (cold start), resolve slot placeholders from project structure:

| Placeholder | Resolution Strategy |
|-------------|-------------------|
| `<START_CMD>` | From `package.json` `scripts.dev`/`scripts.start` (Node), `go run` + entry point from `cmd/*/main.go` (Go), `cargo run` (Rust), `python -m <module>` from `pyproject.toml` (Python) |
| `<PORT>` | From Convention if available; else common defaults: Node 3000/5173 (Vite), Go 8080, Python 8000 (Django/FastAPI), Rust 3000. Check `package.json` scripts for `--port`/`-p` flags. |
| `<HEALTH_URL>` | `http://localhost:<PORT>/healthz` (default). Convention may specify a different endpoint (e.g., `/api/health`, `/ready`). |
| `<PID_FILE>` | `.forge/<key>.pid` (per-surface PID isolation in `.forge/` directory). |

**`# user-customized` marker scope**:
- **Lifecycle recipes** (dev, probe, test, teardown, test-setup): MUST include `# user-customized` — users may customize server commands, probe endpoints, test runners.
- **Gate recipes** (compile, fmt, lint, unit-test): MUST include `# user-customized` — users may customize tool flags, file scopes.
- **Aggregate recipes** (`<prefix>` with no verb): Do NOT include `# user-customized` — these are auto-generated orchestration that calls sub-recipes; users should customize the sub-recipes instead.

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
<prefix>test-setup:   (mobile only)
<prefix>dev:          (if applicable)
<prefix>probe:        (if applicable)
<prefix>test:
<prefix>teardown:
<prefix>:              (aggregate, if applicable)

... (repeat for each surface)

# --- end forge standard recipes ---
```

Use the surface **key** for the group name. For named surfaces, use the key (e.g., `[group: backend]`). For scalar surfaces (no key), use the type (e.g., `[group: tui]`).

### Step 4: Verify and Self-Correct

Three-phase verification: (1) consistency check against surface rule contracts, (2) dry-run catches syntax/structure errors, (3) actual execution catches runtime errors.

#### 4a. Phase 1 — Consistency Verification

After generating the justfile, verify that the generated recipes match the Recipe Invocation Contract defined in each surface rule file. For each surface in `SURFACES_LIST`:

1. Read the surface rule file's **Recipe Invocation Contract** table. Note: contract tables use the surface type as prefix (e.g., `api-dev`) for illustration — for named surfaces, the actual recipe name uses the surface key (e.g., `backend-dev` for `backend=api`); for scalar surfaces, the prefix is omitted (e.g., `dev`). Apply this mapping when matching.
2. For each recipe listed in the contract, verify the generated justfile contains a recipe with:
   - **Matching name**: The recipe name in the justfile matches the contract after applying the prefix mapping (e.g., contract says `api-dev`, actual recipe is `backend-dev` for `backend=api` — this is a match).
   - **Matching signature**: The recipe's parameter signature matches the contract (e.g., `api-test journey=''` matches `just api-test [journey]`).
   - **Matching platform variants**: Both `[linux]` and `[windows]` variants exist.
3. Report any mismatches:
   - Missing recipe: "CONTRACT MISMATCH: `<recipe-name>` defined in surface rule contract but not found in justfile"
   - Signature mismatch: "CONTRACT MISMATCH: `<recipe-name>` signature differs — contract: `<expected>`, actual: `<found>`"
   - Missing platform variant: "CONTRACT MISMATCH: `<recipe-name>` missing `[linux]` or `[windows]` variant"

If any mismatches are found, fix the justfile before proceeding to Phase 2.

This consistency check ensures that `init-justfile` (producer) and `run-tests` (consumer) share a single source of truth for recipe contracts defined in surface rule files.

#### 4b. Phase 2 — Dry-run (syntax check)

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
# Note: ci recipe is excluded — it aggregates other recipes (lint + compile + unit-test), which are individually verified above.
```

For surface recipes, verify each generated recipe (`<prefix>` is `<key>-` for named surfaces or empty for scalar):

```bash
just --dry-run <prefix>test-setup # (mobile only)
just --dry-run <prefix>dev       # (if applicable)
just --dry-run <prefix>probe     # (if applicable)
just --dry-run <prefix>test
just --dry-run <prefix>teardown
```

Fix any syntax failures before proceeding to Phase 3.

#### 4c. Phase 3 — Actual execution (runtime check)

Execute each recipe for real to catch runtime errors. Recipes are classified by execution safety:

| Category                                              | Recipes                    | Method                                                                                    |
| ----------------------------------------------------- | -------------------------- | ----------------------------------------------------------------------------------------- |
| **Safe** (fast, no side effects)                      | `compile`, `lint`, `check` | Execute directly                                                                          |
| **Destructive** (modifies files or creates artifacts) | `build`, `fmt`, `clean`    | Execute directly (artifacts can be cleaned; fmt changes are welcome)                      |
| **Idempotent** (installs dependencies)                | `install`                  | Execute directly                                                                          |
| **Long-running** (starts servers)                     | `<prefix>dev`                | Execute with timeout (10s), kill after timeout -- success = process still alive at timeout |
| **Expensive** (runs full test suite)                  | `unit-test`, `<prefix>test`  | Skip actual execution; verified by `--dry-run` only                                       |

For long-running recipes (`<prefix>dev`): execute via `timeout 10 just <prefix>dev 2>&1 || true`. A crash before timeout ("missing script", "can't load package") is a runtime failure.

#### 4d. Self-correction rules

When a recipe fails in Phase 3, analyze the error and apply corrections per `rules/self-correction.md`:

1. Match the error against known error patterns (npm missing scripts, Go package issues, missing linters/formatters).
2. Edit the justfile to apply the correction.
3. Re-run the failed command (actual execution, same method as Phase 3).
4. If it still fails after 2 attempts, leave the recipe as-is and report the failure in the output.

#### 4e. Report verification results

Output one line per recipe: `[ok]` or `[fix]` + recipe name + resolved command + method (executed / dry-run only / self-corrected). Summarize auto-corrections at the end.

### Step 5: Output Confirmation

Output format (adapt per project):

```
Created justfile with <surface-aware | standard> forge targets (<Language> project)

Surfaces:              (omit if no surfaces)
  <key> (<type>): <generated recipes>

Language targets:
  just <target>        -> <resolved command>
  ...

Surface targets:       (omit if no surfaces)
  just <prefix><verb>  -> <resolved command>
  ...

Convention: <path> (<framework>) | No Convention file found. Run `/forge:test-guide` to create.
Edit justfile to customize commands. Recipes marked `# user-customized` will be preserved on re-generation.
```

## Notes

- **just >= 1.50.0**: supports `[arg]` named option syntax and `[linux]`/`[windows]` platform attributes; surface recipes use dual-platform variants.
- **Breaking change — surface prerequisite**: Projects without surface configuration must run `forge init` first (Outcome B in Step 0). This is a deliberate change from the previous behavior (which generated a basic justfile without surface recipes). Upgrade path: run `forge init` to configure surfaces, then re-run `/init-justfile`.
- **Two-layer model**: `unit-test` is language-level (fast, per-task submit gate); `<prefix>test` is surface-level (functional tests for cli/tui/api, e2e tests for web/mobile). Forge is surface-agnostic -- it calls `just <prefix>test` based on task surface-key. Test type terminology follows the inline model below.
- **Test directory path**: Single-surface projects use `tests/<journey>/` (no surface key layer). Multi-surface projects use `tests/<key>/<journey>/` (key-based directory layer). The `<prefix>test` recipe body must resolve scripts from the correct path.

<!-- INLINE:origin=test-guide/references/test-type-model.md -->
**Surface -> Test Type mapping**:

| Surface | Test Type | Verification | Execution Model |
|---------|-----------|-------------|-----------------|
| `cli` | CLI Functional Test | Exit code + stdout + stderr | Subprocess |
| `tui` | Terminal Functional Test | Terminal output + stdin interaction | Subprocess + stdin pipe |
| `api` | API Functional Test | HTTP status + response body + headers | HTTP client |
| `web` | Web E2E Test | DOM visibility + user interaction + URL change | Browser automation |
| `mobile` | Mobile E2E Test | UI visibility + user interaction + screen ID | Maestro YAML / manual |

Key: "Functional" = protocol-level call verification (subprocess/HTTP); "E2E" = device-level automation (browser/mobile). The "e2e" label applies **only** to web and mobile surfaces.
<!-- END INLINE:origin=test-guide/references/test-type-model.md -->
- **Recipe naming**: Scalar surfaces (no key) produce prefix-less recipes (`test`, `dev`, `teardown`). Named surfaces produce `<key>-` prefixed recipes (e.g., `admin-panel-dev`, `payment-api-test`). Multi-surface projects always use named keys, so each surface gets its own prefix.
- **Targets invoked by forge skills**: `compile`, `unit-test`, `<prefix>test`, `<prefix>teardown`, `install`. The remaining targets are for manual use.
- **Cold start**: When no Convention files exist, the agent generates recipes from its built-in knowledge of common patterns for the detected language. These recipes use conservative defaults and may need manual adjustment.
- **Agent-driven generation**: No language templates are used. The agent detects language via marker files, reads recipe contracts from surface rules, acquires framework knowledge from Convention files, and generates concrete recipe commands. Server lifecycle patterns are referenced from `rules/server-lifecycle.md`.

<EXTREMELY-IMPORTANT>
- MANUAL-ONLY. Do NOT auto-invoke this skill from other skills or agents. Only invoke when user explicitly runs `/init-justfile`.
- If an existing justfile lacks forge boundary markers and `--force` is not set, you MUST prompt the user before overwriting. Never silently destroy user customizations.
- Only the section between `# --- forge standard recipes ---` / `# --- end forge standard recipes ---` markers may be replaced. Recipes outside markers must be preserved verbatim.
- After writing, you MUST run the verification steps (consistency check + dry-run + actual execution) and report all results.
- CLI/TUI surfaces MUST NOT generate dev, probe, or aggregate recipes.
- `# user-customized` marked recipes MUST be preserved during re-generation.
- No language templates are used. The agent generates recipes directly from detected language + Convention knowledge + surface rule contracts.
- Surface data source is `forge surfaces` (text mode). Scalar surfaces (no `=` in output line) produce prefix-less recipes (`test`, `dev`, `teardown`). Named surfaces produce `<key>-` prefixed recipes.
- Server lifecycle patterns MUST be sourced from `rules/server-lifecycle.md` — do not generate PID tracking, idempotent start, or health check code from scratch.
</EXTREMELY-IMPORTANT>
