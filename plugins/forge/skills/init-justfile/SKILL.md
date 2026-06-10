---
name: init-justfile
description: Scaffold a Justfile with standard forge targets for the current project, with surface-aware recipe generation.
allowed-tools: Bash Read Write Edit
disable-model-invocation: true
argument-hint: '[--force]'
---

# /init-justfile

MANUAL-ONLY. Do NOT auto-invoke — only when user explicitly asks to `/init-justfile`.

Generate a Justfile with standard forge targets. When surfaces are configured, call `forge justfile scaffold` to produce surface-aware recipes (dev/probe/test/teardown) differentiated by surface type.

## Prerequisites

**Install just (>= 1.50.0)**: `brew install just` (macOS/Linux), `scoop install just` (Windows), or `cargo install just`. Requires >= 1.50.0 for `[arg]` named option syntax and `[linux]`/`[windows]` attributes.

## Parameters

| Parameter | Values  | Default | Description                                      |
| --------- | ------- | ------- | ------------------------------------------------ |
| `--force` | (flag)  | false   | Overwrite existing justfile without confirmation |

## Recipe Naming Model

`<prefix>` is `<key>-` for named surfaces (e.g., `app-test`) or empty for scalar surfaces (e.g., `test`).

| Surface Type | Lifecycle Recipes | Quality Recipes | Aggregate |
| ------------ | ----------------- | --------------- | --------- |
| `cli`        | `test`, `teardown` | `compile`, `fmt`, `lint`, `unit-test` | 无 |
| `tui`        | `test`, `teardown` | `compile`, `fmt`, `lint`, `unit-test` | 无 |
| `api`/`web`  | `dev`, `probe`, `test`, `teardown`, `<key>` | `compile`, `fmt`, `lint`, `unit-test` | `<key>` (dev->probe->test->teardown) |
| `mobile`     | `test-setup`, `dev`, `probe`, `test`, `teardown`, `<key>` | `compile`, `fmt`, `lint`, `unit-test` | `<key>` (test-setup->dev->probe->test->teardown) |

**Test directory path**: Single-surface projects use `tests/<journey>/`. Multi-surface projects use `tests/<key>/<journey>/`.

**Recipe signatures**: `unit-test` (no params), `<prefix>test [journey]` (optional filter), `<prefix>probe` (no params).

## Process Flow

```
Step 0: forge surfaces -> SURFACES_LIST (key + type + form)
Step 1: Detect languages per surface + load Convention -> record slot values
Step 2: Check existing justfile (boundary markers + user-customized protection)
Step 3: forge justfile scaffold per surface -> fill <<PLACEHOLDER>> -> --aggregate -> boundary marker merge
Step 4: Verify (just --list + dry-run + actual execution + self-correction)
Step 5: Output confirmation
```

### Step 0: Detect Surfaces

```bash
forge surfaces 2>/dev/null
```

Parse each line: `key=type` (named) or `type` alone (scalar). Collect as `SURFACES_LIST`. If empty/failed, output "No surfaces configured. Run `forge init` first." and abort.

### Step 1: Detect Languages and Load Convention

**1a. Language detection**: Scan for marker files (`go.mod`=Go, `package.json`=Node/TS, `Cargo.toml`=Rust, `pyproject.toml`/`setup.py`=Python, `pom.xml`/`build.gradle`=Java, `build.gradle.kts`=Kotlin, `Gemfile`=Ruby, `*.csproj`/`*.sln`=C#/.NET). Read marker to extract version, dependency manager, tooling.

**1b. Per-surface detection**: Each surface in `SURFACES_LIST` scans its working directory independently using 1a's table. Multi-surface projects may have different languages per surface.

**1c. Convention loading**: Read `docs/conventions/testing/<surface>/core.md` for each surface. Extract: framework, file pattern, test runner, build tags, result format flags. Legacy fallback: flat files with `domains` frontmatter. Cold start: use agent built-in knowledge.

**1d. Cold start fallback**: When Convention absent: (1) infer from project files (`package.json` scripts, `Makefile`, go module), (2) language defaults (Node: 3000, Go: 8080), (3) leave unresolvable as `<<PLACEHOLDER>>`, list in report.

### Step 2: Check Existing Justfile

- No justfile: proceed to Step 3.
- Boundary markers present (`# --- forge standard recipes ---` / `# --- end forge standard recipes ---`): check completeness via `just --list`, skip to Step 4 if complete.
- No markers + `--force`: proceed to Step 3.
- No markers + no `--force`: prompt "Overwrite? (y/n)", abort if declined.

<HARD-RULE>
If an existing justfile lacks forge boundary markers and `--force` is not set, you MUST prompt the user before overwriting. Never silently destroy user customizations.
</HARD-RULE>

### Step 3: Generate Recipes and Assemble Justfile

**3a. CLI scaffold**: For each surface, run `forge justfile scaffold --type <type> --key <key>` (named) or `--type <type>` (scalar). Fill `<<PLACEHOLDER>>` slots from Convention (priority) > project inference > language defaults. Then `forge justfile scaffold --aggregate` for `install`/`ci`/`clean`.

**3b. Boundary marker merge**: Replace content between markers (inclusive). No markers: write full new file.

<HARD-RULE>
Only the section between `# --- forge standard recipes ---` / `# --- end forge standard recipes ---` markers may be replaced. Recipes outside markers MUST be preserved verbatim. `# user-customized` marked recipes MUST be preserved during re-generation.
</HARD-RULE>

**3c. Organization**: Group recipes as `[group: language]` (compile/build/install/clean/ci), `[group: language-test]` (unit-test/lint/fmt/check), `[group: <surface-key>]` (surface lifecycle + aggregate). Use surface key for group name; scalar uses type.

### Step 4: Verify and Self-Correct

**4a. Completeness**: `just --list` — confirm all recipes parseable. Fix errors before proceeding.

**4b. Dry-run**: `just --dry-run` each recipe to verify syntax.

**4c. Actual execution**: Safe recipes (compile/lint/check) execute directly. Destructive (build/fmt/clean) execute directly. Idempotent (install) execute directly. Long-running (`<prefix>dev`) use `timeout 10s`. Expensive (unit-test/`<prefix>test`) skip, dry-run only.

**4d. Self-correction**: On failure, apply fixes per `rules/self-correction.md`. Max 2 retries per recipe.

**4e. Report**: One line per recipe: `[ok]` or `[fix]` + name + command + method.

### Step 5: Output Confirmation

```
Created justfile with <surface-aware | standard> forge targets (<Language> project)

Surfaces:              (omit if no surfaces)
  <key> (<type>): <generated recipes>

Language targets:
  just <target>        -> <resolved command>

Surface targets:       (omit if no surfaces)
  just <prefix><verb>  -> <resolved command>

Convention: <path> (<framework>) | No Convention file found. Run `/forge:test-guide` to create.
Edit justfile to customize commands. Recipes marked `# user-customized` will be preserved on re-generation.
```

## Notes

- **Two-layer model**: `unit-test` = language-level (per-task submit gate); `<prefix>test` = surface-level (functional/e2e).
- **Recipe generation**: MUST use `forge justfile scaffold` CLI. Do NOT generate PID/lifecycle code from scratch.
- **Targets invoked by forge skills**: `compile`, `unit-test`, `<prefix>test`, `<prefix>teardown`, `install`.

<!-- INLINE from test-guide/references/test-type-model.md @ v3.0.0-rc.53 -->
**Surface -> Test Type mapping**:

| Surface | Test Type | Verification | Execution Model |
|---------|-----------|-------------|-----------------|
| `cli` | CLI Functional Test | Exit code + stdout + stderr | Subprocess |
| `tui` | Terminal Functional Test | Terminal output + stdin interaction | Subprocess + stdin pipe |
| `api` | API Functional Test | HTTP status + response body + headers | HTTP client |
| `web` | Web E2E Test | DOM visibility + user interaction + URL change | Browser automation |
| `mobile` | Mobile E2E Test | UI visibility + user interaction + screen ID | Maestro YAML / manual |
<!-- END INLINE:origin=test-guide/references/test-type-model.md -->

<EXTREMELY-IMPORTANT>
- MANUAL-ONLY. Do NOT auto-invoke this skill.
- Boundary marker + `--force` protection: never silently overwrite user customizations.
- `# user-customized` recipes MUST be preserved during re-generation.
- CLI/TUI surfaces MUST NOT generate dev, probe, or aggregate recipes.
- Recipe generation MUST use `forge justfile scaffold` CLI command.
</EXTREMELY-IMPORTANT>
