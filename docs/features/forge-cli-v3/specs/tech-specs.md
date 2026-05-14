---
feature: "forge-cli-v3"
generated: "2026-05-14"
status: draft
---

# Technical Specifications: Forge CLI v3

## Command Structure

### TECH-001: Cobra Command Grouping Convention

**Requirement**: Commands MUST be organized into domain groups using Cobra's `GroupID` mechanism. Each group corresponds to a business domain: task, e2e, forensic, profile, prompt. Top-level commands (feature, probe, cleanup, quality-gate, verify-task-done, version) are NOT grouped.
**Scope**: [LOCAL]
**Source**: prd-spec.md > Functional Specs > Command Structure

### TECH-002: Binary and Module Naming

**Requirement**: Go module MUST be named `forge-cli` (replacing `task-cli`). Binary output MUST be `forge` (replacing `task`). Directory MUST be `forge-cli/` (replacing `task-cli/`). The `pkg/version/version.go` Name constant MUST be updated.
**Scope**: [LOCAL]
**Source**: prd-spec.md > Scope > In Scope

### TECH-003: Command Renaming Mapping

**Requirement**: The following command renames MUST be applied exactly:

| Old | New | Rationale |
|-----|-----|-----------|
| `task` | `forge` | Brand unification |
| `task record` | `forge task submit` | Disambiguate noun/verb |
| `task check` | `forge task check-deps` | Clarify check target |
| `task validate` | `forge task validate-index` | Clarify validation target |
| `task verify-completion` | `forge task verify-task-done` | Clarify verification target |
| `task all-completed` | `forge quality-gate` | Reflect actual behavior |
| `task prompt` | `forge prompt get-by-task-id` | Disambiguate for AI context |
| `task template` | (DELETE) | No consumers |

**Scope**: [LOCAL]
**Source**: prd-spec.md > Functional Specs > Naming Changes

## Profile System

### TECH-004: Profile-Aware E2E Architecture

**Requirement**: E2E commands MUST read the `profile` field from `.forge/config.yaml` to determine which test suite to execute. Profile detection scans project structure for framework-specific config files (e.g., `playwright.config.ts` -> web-playwright). Supported profiles are defined in a registry; unknown profile values MUST be rejected.
**Scope**: [CROSS]
**Source**: prd-spec.md > User Stories > Story 5 & Story 8

## Performance

### TECH-005: Startup Latency Budget

**Requirement**: `forge --help` response time MUST be <= baseline + 50ms, where baseline is the median of 3 runs of `task --help` before migration. Binary size increase MUST be <= 500KB. Build time increase MUST be <= 10 seconds.
**Scope**: [LOCAL]
**Source**: prd-spec.md > Other Notes > Performance Requirements

## Error Handling

### TECH-006: Stderr-Only Error Output Pattern

**Requirement**: All error diagnostics MUST be written to stderr, never to stdout. stdout is reserved for command output data. Error messages MUST follow the format: `<context>: <specific-detail>` (e.g., "task not found: T-impl-1", "unknown profile: bad-value").
**Scope**: [CROSS]
**Source**: prd-spec.md > Error Handling > Command-level failures (pattern across all commands)

### TECH-007: Behavior Equivalence Guarantee

**Requirement**: All migrated commands MUST produce identical exit codes (0/1/2) and stdout format compared to the pre-migration `task` CLI. Only `--help` output is allowed to change due to the new grouping structure.
**Scope**: [LOCAL]
**Source**: prd-spec.md > Other Notes > Performance Requirements

## Data Integrity

### TECH-008: No Data Migration Required

**Requirement**: The `.forge/` directory structure, state.json format, and config.yaml format MUST remain unchanged. No data migration scripts are needed.
**Scope**: [LOCAL]
**Source**: prd-spec.md > Other Notes > Data Requirements

### TECH-009: File Lock Mechanism Preservation

**Requirement**: The index.json.lock file-based locking mechanism MUST be preserved unchanged during the rename. Concurrent safety MUST be verified post-migration.
**Scope**: [LOCAL]
**Source**: prd-spec.md > Other Notes > Security Requirements
