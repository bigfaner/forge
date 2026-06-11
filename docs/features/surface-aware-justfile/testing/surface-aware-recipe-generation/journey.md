---
feature: "surface-aware-justfile"
journey: "surface-aware-recipe-generation"
risk_level: "Medium"
sources:
  - docs/features/surface-aware-justfile/prd/prd-user-stories.md
  - docs/features/surface-aware-justfile/prd/prd-spec.md
generated: "2026-05-26"
---

# Journey: Surface-Aware Recipe Generation

**Risk Level**: Medium

<!-- Risk Classification Criteria:
  Medium = Workflow involves multi-step interaction without irreversible side effects.
  init-justfile generates/overwrites justfile recipes but has user-customized protection.
  No data loss risk -- user-customized recipes are preserved.
-->

## Overview

A Forge user configures the `surfaces` field in `.forge/config.yaml` and runs `init-justfile` to automatically generate surface-type-specific dev/test/probe/teardown recipes, so that web/api projects get the correct "start -> wait -> test -> cleanup" orchestration and cli/tui projects get "build -> test" orchestration.

## Setup

- `.forge/config.yaml` exists with a `surfaces` field defining at least one surface entry (e.g., `admin-panel: web`)
- The project has a recognized language (Go, Node.js, etc.) so that base recipes (compile/build/lint/fmt) can be generated
- The `just` binary is installed and version >= 1.4.0

## Happy Path

### Step 1: Configure surfaces in config.yaml

**User Action**: Define surface entries in `.forge/config.yaml`, e.g., `surfaces: {admin-panel: web}`

**Expected Result**: The config.yaml `surfaces` field contains a valid map of surface-key to surface-type (one of web/api/cli/tui/mobile).

### Step 2: Run init-justfile

**User Action**: Execute `init-justfile` skill via Forge CLI

**Expected Result**: init-justfile detects the surface type from config.yaml, loads the corresponding surface rule file (`skills/init-justfile/rules/surfaces/web.md`), and generates surface-specific recipes.

### Step 3: Verify generated justfile contains surface-specific recipes

**User Action**: Inspect the generated justfile

**Expected Result**: The justfile contains:
- For web/api surfaces: `dev` (background start), `probe` (retry polling), `test`, `test-teardown` recipes
- For cli/tui surfaces: `dev`, `test` recipes only (no `run`, no `probe`)
- Each recipe includes `[linux]`/`[windows]` dual-platform variants

### Step 4: Verify user-customized protection

**User Action**: Manually modify a generated recipe in the justfile, then re-run init-justfile

**Expected Result**: Recipes marked with `# user-customized` comment are preserved unchanged. init-justfile outputs a diff summary showing what would have changed. To overwrite, the user must pass `--force-regenerate`.

### Step 5: Verify mixed-project recipe generation

**User Action**: Configure multiple surfaces (e.g., `{admin-panel: web, payment-service: api}`) and run init-justfile

**Expected Result**: The generated justfile contains:
- Per-surface-key prefixed recipes: `dev-admin-panel`, `dev-payment-service`, `probe-admin-panel`, `probe-payment-service`, etc.
- Aggregation recipes: `dev` (starts all services in dependency order), `test` (runs all test sequences)
- An orchestration order comment at the top of the justfile

## Edge Cases

### Step 1b: Invalid surface type in config.yaml

**Precondition**: config.yaml `surfaces` field contains an unrecognized surface type (e.g., `{my-surface: desktop}`)

**User Action**: Run init-justfile

**Expected Result**: init-justfile outputs a descriptive error to stderr indicating the unsupported surface type and lists the 5 supported types (web/api/cli/tui/mobile). Exit with non-zero code.

### Step 2b: No surfaces configured in config.yaml

**Precondition**: `.forge/config.yaml` exists but has no `surfaces` field, or the field is empty

**User Action**: Run init-justfile

**Expected Result**: init-justfile generates only language-template-based recipes (compile/build/lint/fmt) with no orchestration recipes, producing output identical to the current (pre-feature) behavior. Zero regression.

### Step 3b: just version below 1.4.0

**Precondition**: The installed `just` binary version is below 1.4.0

**User Action**: Run init-justfile

**Expected Result**: init-justfile outputs an error to stderr with the current version and the required version ("just >= 1.4.0"), then exits with exit code 2 (blocking).

### Step 4b: Surface rule file missing

**Precondition**: A supported surface type is configured (e.g., `web`) but the corresponding rule file (`skills/init-justfile/rules/surfaces/web.md`) does not exist

**User Action**: Run init-justfile

**Expected Result**: init-justfile outputs an error to stderr with the missing file path and a recovery hint ("run init-justfile to regenerate rule files"), then exits with exit code 2.

### Step 5b: Malformed surfaces config (not map<string,string>)

**Precondition**: config.yaml `surfaces` field has invalid format (e.g., a scalar value or a list instead of a map)

**User Action**: Run init-justfile

**Expected Result**: init-justfile outputs a descriptive YAML parse error to stderr with a recovery hint ("check .forge/config.yaml surfaces field format, should be map<string, string>, e.g., {admin-panel: web}"), then exits with exit code 2.

### Step 6b: Surface-key contains invalid characters

**Precondition**: A surface-key in config.yaml contains characters not allowed in just recipe names (e.g., `my/surface` or `surface+key`)

**User Action**: Run init-justfile

**Expected Result**: init-justfile validates surface-key names against `[a-zA-Z0-9_-]` and outputs an error for any invalid key, suggesting valid alternatives.

## Journey Invariants

- init-justfile never silently overwrites a recipe marked with `# user-customized` without `--force-regenerate`
- Projects without `surfaces` configuration always produce output identical to the pre-feature behavior (zero regression)
- All generated recipes include dual-platform (`[linux]`/`[windows]`) variants where applicable
- cli/tui surfaces never generate `run` or `probe` recipes
- Mixed-project aggregation recipes always list services in dependency order (api before web)
