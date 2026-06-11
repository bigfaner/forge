---
feature: "surface-aware-justfile"
journey: "surface-key-migration"
risk_level: "High"
sources:
  - docs/features/surface-aware-justfile/prd/prd-user-stories.md
  - docs/features/surface-aware-justfile/prd/prd-spec.md
generated: "2026-05-26"
---

# Journey: Surface-Key Migration and Data Model Extension

**Risk Level**: High

<!-- Risk Classification Criteria:
  High = Workflow involves state mutation, data loss risk, or irreversible operations.
  This journey covers Go struct field migration (Scope->SurfaceKey), resolveScope() rewrite,
  16 prompt template variable changes, and dead code removal. Data model changes are
  irreversible at the code level and affect 7+ components.
-->

## Overview

A Forge plugin developer migrates the surface-key value domain from a fixed enumeration (frontend/backend) to user-defined surface-key names (while surface-type remains 5 fixed types), extends the Task data model with surface-key and surface-type fields, and updates all downstream components (breakdown-tasks, quick-tasks, forge task add, quality-gate fix-task, init-justfile, prompt templates) to use the new unified identifiers.

## Setup

- The project has `.forge/config.yaml` with a `surfaces` field defining surface-key to surface-type mappings (e.g., `{admin-panel: web, payment-service: api}`)
- The `forge surfaces` CLI command is available and functional (prerequisite)
- The existing codebase uses `scope` field with hardcoded `frontend`/`backend` values in various places

## Happy Path

### Step 1: Run forge surfaces CLI to verify surface detection

**User Action**: Execute `forge surfaces <path>` for a configured project path

**Expected Result**: The CLI returns a surface-key and surface-type via longest-prefix-match. For example, `forge surfaces frontend/src` returns `{"surface-key": "admin-panel", "surface-type": "web"}`.

### Step 2: Verify Task Go struct migration

**User Action**: After implementation, inspect `task/types.go`

**Expected Result**: The `Scope` field is replaced by `SurfaceKey` (string). A new `SurfaceType` field (string) is added. `GetSurfaceKey()` provides backward-compatible access. The JSON serialization includes `"surfaceKey"` and `"surfaceType"` fields.

### Step 3: Verify resolveScope() rewrite

**User Action**: Inspect `prompt.go` resolveScope() function

**Expected Result**: resolveScope() performs a surfaces map collection query instead of hardcoded projectType matching. It uses `forge surfaces` CLI output to determine surface-key dynamically.

### Step 4: Verify breakdown-tasks generates tasks with surface fields

**User Action**: Run breakdown-tasks for a project with configured surfaces

**Expected Result**: Generated task files have `surface-key` and `surface-type` in their frontmatter. Values match the configured surfaces (e.g., `surface-key: admin-panel`, `surface-type: web`).

### Step 5: Verify forge task add inherits surface fields

**User Action**: Run `forge task add` with a source task that has surface-key/surface-type

**Expected Result**: The new task inherits `surface-key` and `surface-type` from the source task. When no source task exists and the project has a single surface, the unique surface-type is auto-filled.

### Step 6: Verify quality-gate fix-task infers surface from file path

**User Action**: A quality-gate fix-task is created from a failing test file path

**Expected Result**: The fix-task's surface-key and surface-type are inferred from the failing file path using `forge surfaces <path>` longest-prefix-match.

### Step 7: Verify zero-regression for projects without surfaces

**User Action**: Run the full workflow on a project with no `surfaces` configuration in config.yaml

**Expected Result**: All behavior is identical to the pre-feature baseline. No surface-key or surface-type fields appear in generated tasks. resolveScope() falls back gracefully.

## Edge Cases

### Step 1b: forge surfaces CLI returns no match

**Precondition**: The queried path does not match any configured surface entry

**User Action**: Run `forge surfaces unknown-dir`

**Expected Result**: CLI exits with code 1. stderr contains an error message with recovery hint ("run forge init to configure surfaces"). Downstream components handle this gracefully.

### Step 2b: Old task files with `scope: frontend` exist

**Precondition**: Existing task files have `scope: frontend` or `scope: backend` in their frontmatter

**User Action**: Attempt to read the task via `forge task status` or similar

**Expected Result**: `forge task migrate` can automatically migrate `scope` to `surface-key` + `surface-type`. Before migration, task read commands return a blocking error (exit 2) to prevent silent data loss.

### Step 3b: forge surfaces CLI execution fails (not installed or wrong version)

**Precondition**: `forge surfaces` command is not available (CLI not installed or version too old)

**User Action**: Run breakdown-tasks or any component that calls forge surfaces

**Expected Result**: Component outputs error to stderr with the CLI output and a recovery hint ("check forge CLI is installed and version >= required version"). Exits with exit code 1 (retryable).

### Step 4b: Multiple surfaces match the same file path (ambiguous)

**Precondition**: A file path matches multiple surface entries via longest-prefix-match (e.g., overlapping path prefixes in config)

**User Action**: Run `forge surfaces <ambiguous-path>`

**Expected Result**: The CLI selects the longest (most specific) prefix match. If two entries have identical prefix length, an error is returned indicating ambiguous configuration.

### Step 5b: forge task add with no source task and multi-surface project

**Precondition**: `forge task add` is called without a source task, and the project has multiple surfaces configured

**User Action**: Run `forge task add` without `--source-task-id`

**Expected Result**: Surface-type cannot be auto-filled (ambiguous). The command requires explicit `--surface-type` flag, or outputs an error listing available surfaces.

### Step 6b: 16 prompt templates SURFACE_KEY variable mismatch

**Precondition**: Some prompt templates still reference the old `frontend`/`backend` values

**User Action**: Run any prompt-based operation that uses these templates

**Expected Result**: All 16 prompt templates use the new `SURFACE_KEY` variable with user-defined surface-key values at runtime. No hardcoded `frontend`/`backend` references remain in template logic.

### Step 7b: Dead code remnants (extractTestTypeArg, genScriptBases)

**Precondition**: After migration, dead code from the old scope system remains

**User Action**: Compile the project

**Expected Result**: `extractTestTypeArg()` and `genScriptBases` functions are removed. No compilation errors. All callers updated to use the new surface-based APIs.

### Step 8b: Surface-key in just recipe name has invalid characters

**Precondition**: A user-defined surface-key contains characters incompatible with just recipe names

**User Action**: Run init-justfile with such a surface-key

**Expected Result**: init-justfile validates surface-key against `[a-zA-Z0-9_-]` before generating recipes. Invalid characters trigger a descriptive error.

## Journey Invariants

- `surface-type` always belongs to the fixed set {web, api, cli, tui, mobile} -- never user-defined
- `surface-key` is always user-defined and unique within a project's config.yaml
- Migration is phased: Phase 1 (data model) -> Phase 2 (upstream adapters) -> Phase 3 (downstream consumers) -- strict sequential dependency
- Projects without `surfaces` configuration produce identical output to the pre-feature baseline (zero regression guarantee)
- `forge task migrate` must exist before any task read operations work on old-format task files (blocking error prevents silent data loss)
- All 7+ components surface-key value domains are synchronized -- no component uses hardcoded frontend/backend after migration
