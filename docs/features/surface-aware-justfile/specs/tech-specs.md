---
feature: "surface-aware-justfile"
generated: "2026-05-26"
status: draft
---

# Technical Specifications: Surface-Aware Justfile

## CLI Interface

### TECH-001: forge surfaces --json Flag Specification

**Requirement**: The `forge surfaces` CLI command supports a `--json` bool flag that outputs structured JSON for machine consumption by skills. Three sub-modes:
- **List mode** (no path, no --types): `{"surfaces": [{"key": "admin-panel", "type": "web"}, ...]}`
- **Query mode** (with path): `[{"key": "admin-panel", "type": "web"}]`; no match returns `[]` with exit 0
- **Types mode** (--types flag): `{"types": ["api", "web"]}`

When surfaces config is missing (all modes): stderr `{"error": "no surface configured; run 'forge init' to configure surfaces"}`, exit 1.
**Scope**: [CROSS]
**Source**: tech-design.md "Interface 1b: CLI --json flag specification"

### TECH-002: --json Mode stderr Format Override

**Requirement**: When `--json` flag is active, ALL output (stdout AND stderr) MUST use structured JSON format. This is an explicit exception to the plain-text stderr format convention (`<context>: <specific-detail>`). Rationale: `--json` consumers are machines (skills via Bash tool), requiring stderr to be parseable by JSON parsers without ambiguity. All error paths under `--json` must use `json.NewEncoder(cmd.ErrOrStderr()).Encode()` with `{"error": "..."}` format; `fmt.Fprintf(os.Stderr, ...)` is prohibited.
**Scope**: [CROSS]
**Source**: tech-design.md "Interface 1b: --json stderr format override"

## Surface Rule Files

### TECH-003: Surface Rule File Format Convention

**Requirement**: Surface rule files follow the path pattern `rules/surfaces/<type>.md` and are consumed by both init-justfile (recipe generation) and run-tests (orchestration). Each file contains: (1) orchestration sequence table with exit code handling per step, (2) recipe invocation contract table defining just signature, exit code semantics, and (3) journey filter strategy mapping @tag to surface type. Files are loaded dynamically based on surface-type detected at runtime.
**Scope**: [CROSS]
**Source**: tech-design.md "Interface 2: Surface rule file format"

### TECH-004: Recipe Naming Convention for Mixed Projects

**Requirement**: In mixed projects (multiple surfaces), recipes use `<action>-<surface-key>` naming pattern (e.g., `dev-admin-panel`, `probe-payment-service`). Aggregate recipes without prefix serve as default entry points (e.g., `dev` starts all dev servers in dependency order). Surface-key in recipe names must use `[a-zA-Z0-9_-]` characters only. Orchestration order is recorded in justfile header comments for run-tests parsing.
**Scope**: [CROSS]
**Source**: prd-spec.md "混合项目生成与编排流程"

## Data Model

### TECH-005: Task Struct Surface Fields

**Requirement**: Task Go struct adds `SurfaceKey string` (JSON: `"surface-key,omitempty"`) and `SurfaceType string` (JSON: `"surface-type,omitempty"`). The `Scope` field is removed entirely -- no backward compatibility layer. Old frontmatter with `scope` but no `surface-key` triggers a blocking error (exit 2) via `CheckLegacyScope()` shared function.
**Scope**: [LOCAL]
**Source**: tech-design.md "Model 1: Task struct"

### TECH-006: Cross-Layer Surface Data Propagation

**Requirement**: Surface information flows through a 7-step chain:
1. `.forge/config.yaml` surfaces field -> `forge surfaces` CLI (file read + CLI query)
2. `forge surfaces` CLI -> breakdown-tasks/quick-tasks skill (JSON stdout)
3. Skill -> task frontmatter (YAML fields: surface-key, surface-type)
4. Frontmatter -> index.json Task Go struct (JSON serialization)
5. index.json -> run-tests skill (Go function call)
6. run-tests -> execution strategy rule file (file path: `rules/surfaces/{type}.md`)
7. `forge surfaces` CLI -> init-justfile skill (CLI stdout)

Fallback chain: task frontmatter -> `forge surfaces <path>` longest-prefix-match -> error exit.
**Scope**: [CROSS]
**Source**: prd-spec.md "跨组件数据流"

## Performance

### TECH-007: Surface Rule Loading Performance

**Requirement**: Surface rule file loading MUST NOT add more than 1 second to init-justfile execution time. Just version requirement: >= 1.4.0 (for `[linux]`/`[windows]` recipe attributes). Version check is the first step in both init-justfile and run-tests.
**Scope**: [LOCAL]
**Source**: prd-spec.md "Performance Requirements"

## Prompt Template Migration

### TECH-008: Template Variable Replacement

**Requirement**: 18 prompt templates in `forge-cli/pkg/prompt/data/` and 3 skill templates replace `{{SCOPE}}` with `{{SURFACE_KEY}}`. The `{{TEST_TYPE_ARG}}` variable remains unchanged but its source changes from `extractTestTypeArg()` to direct read of `task.SurfaceType`. Dead code removal: `extractTestTypeArg()` and `genScriptBases()` functions deleted after migration.
**Scope**: [LOCAL]
**Source**: tech-design.md "Prompt template complete list" + "Phase Component Map"
