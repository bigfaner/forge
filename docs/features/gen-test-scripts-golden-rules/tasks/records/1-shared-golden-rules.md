---
status: "completed"
started: "2026-05-21 23:06"
completed: "2026-05-21 23:08"
time_spent: "~2m"
---

# Task Record: 1 Create types/_shared.md cross-type golden rules

## Summary
Created types/_shared.md with 5 cross-type universal golden rules (Isolation, Determinism, Timeout Protection, Idempotency, Resource Cleanup) and 4 shared antipattern guards (Sleep-based waits, Hardcoded configuration, Vacuous assertions, Source-code-level testing). Fully framework-agnostic, structured as Constraint/Rationale/Antipattern guard per principle. Determinism expanded with 3 sub-dimensions. Timeout Protection covers operation-level and test-level scopes with Convention override mechanism.

## Changes

### Files Created
- plugins/forge/skills/gen-test-scripts/types/_shared.md

### Files Modified
无

### Key Decisions
- Each principle structured as Constraint statement + Rationale + Antipattern guard, consistent with the three-layer model: _shared.md (abstract) → type file (type-specific) → Convention (framework)
- Shared antipattern guards extracted from overlapping patterns across all 5 type files: Sleep-based waits (TUI, UI, Mobile), Hardcoded config (API, Mobile), Vacuous assertions (API), Source-code-level testing (CLI, TUI)
- Idempotency principle differentiates between stateful interfaces (API, CLI: create+cleanup) and persistent-side-effect interfaces (UI, TUI, Mobile: repeated interaction must not break subsequent tests)

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] _shared.md defines 5 universal principles: Isolation, Determinism, Timeout Protection, Idempotency, Resource Cleanup
- [x] Each principle has: declarative constraint statement, rationale, and shared antipattern guard
- [x] Determinism principle expanded with sub-dimensions: (a) no random dependency, (b) no external service dependency, (c) no order dependency
- [x] Timeout Protection covers: all I/O operations must have timeout upper bound; default value + Convention override mechanism
- [x] Resource Cleanup covers: tests must not leave behind temp files, background processes, database records, browser sessions
- [x] Shared antipattern guards extracted from duplicates across 5 type files: Sleep-based waits, hardcoded config, vacuous assertions, source-code-level testing
- [x] File is completely framework-agnostic — no language-specific code, commands, or grep patterns
- [x] Frontmatter includes conventions field (empty, as _shared.md is universal)

## Notes
无
