---
status: "completed"
started: "2026-05-19 23:02"
completed: "2026-05-19 23:05"
time_spent: "~3m"
---

# Task Record: 2 Extract phase-detection, db-schema, and existing-code-split rule files

## Summary
Extracted three conditional rule files from SKILL.md: phase-detection.md (three-tier phase detection with explicit/heuristic/fallback tiers, phase-inventory.json format), db-schema.md (schema task creation, AC, breaking classification, scope/dependency rules), and existing-code-split.md (artifact-update + feature sub-task split procedure, thresholds, sub-ID conventions). All files include load conditions, guard clauses, and maintenance notes.

## Changes

### Files Created
- plugins/forge/skills/breakdown-tasks/rules/phase-detection.md
- plugins/forge/skills/breakdown-tasks/rules/db-schema.md
- plugins/forge/skills/breakdown-tasks/rules/existing-code-split.md

### Files Modified
无

### Key Decisions
- Followed the structure of existing ui-placement.md rule file as the canonical format (load condition, guard clause, rules sections, maintenance note)
- Used artifact-relative paths (e.g. design/er-diagram.md, prd/prd-spec.md) consistent with forge distribution model — no hardcoded plugin root paths
- Each file is self-contained and independently understandable without reading SKILL.md

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] rules/phase-detection.md contains three-tier detection (explicit -> heuristic -> fallback)
- [x] rules/phase-detection.md contains explicit detection patterns: flow diagram diamond nodes, PRD sections named Round/Phase/Stage
- [x] rules/phase-detection.md contains heuristic detection patterns with English and Chinese patterns
- [x] rules/phase-detection.md contains fallback: artifact-driven decomposition
- [x] rules/phase-detection.md contains phase-inventory.json format specification
- [x] rules/phase-detection.md has load condition at top
- [x] rules/phase-detection.md has maintenance note listing skeleton dependencies
- [x] rules/db-schema.md contains schema task creation rules (one task per entity)
- [x] rules/db-schema.md contains acceptance criteria: DDL executes, FK references resolve, indexes created
- [x] rules/db-schema.md contains breaking classification: ALTER -> true, CREATE TABLE new -> false
- [x] rules/db-schema.md contains scope assignment: backend
- [x] rules/db-schema.md contains dependency rule: depends on interface tasks
- [x] rules/db-schema.md has load condition at top
- [x] rules/db-schema.md has maintenance note listing skeleton dependencies
- [x] rules/existing-code-split.md contains artifact-update + feature sub-task split procedure
- [x] rules/existing-code-split.md contains sub-ID convention: <seq>.<sub>a and <seq>.<sub>b
- [x] rules/existing-code-split.md contains when-to-apply threshold: >5 files OR cross-layer
- [x] rules/existing-code-split.md contains breaking: true requirement for shared artifact update
- [x] rules/existing-code-split.md contains exclusion for purely additive new code
- [x] rules/existing-code-split.md has load condition at top
- [x] rules/existing-code-split.md has maintenance note listing skeleton dependencies
- [x] All files use skill-relative paths compatible with forge distribution model
- [x] Each file includes guard clause for malformed/empty artifacts
- [x] Each file is independently understandable
- [x] SKILL.md was NOT modified

## Notes
无
