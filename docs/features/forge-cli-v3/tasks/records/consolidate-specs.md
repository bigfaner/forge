---
status: "completed"
started: "2026-05-14 08:25"
completed: "2026-05-14 08:29"
time_spent: "~4m"
---

# Task Record: T-test-5 Consolidate Specs

## Summary
Consolidated specs from forge-cli-v3 PRD: extracted 9 business rules and 9 technical specs. Classified 7 as CROSS (5 biz + 2 tech) and 9 as LOCAL. Integrated CROSS items into project-level dirs: docs/business-rules/task-lifecycle.md (BIZ-001, BIZ-002), docs/business-rules/quality-gate.md (BIZ-004), docs/business-rules/error-reporting.md (BIZ-008, BIZ-009), docs/conventions/profile-system.md (TECH-004), docs/conventions/error-handling.md (TECH-006). No existing project-level files or overlaps found. Preview files and review-choices.md written for traceability. Integration marker written.

## Changes

### Files Created
- docs/features/forge-cli-v3/specs/biz-specs.md
- docs/features/forge-cli-v3/specs/tech-specs.md
- docs/features/forge-cli-v3/specs/review-choices.md
- docs/features/forge-cli-v3/specs/.integrated
- docs/business-rules/task-lifecycle.md
- docs/business-rules/quality-gate.md
- docs/business-rules/error-reporting.md
- docs/conventions/profile-system.md
- docs/conventions/error-handling.md

### Files Modified
无

### Key Decisions
- Classified task state machine and terminal state immutability as CROSS domain invariants since they apply to all features using the task system
- Classified exit code semantics and error message patterns as CROSS conventions since they apply to all CLI commands project-wide
- Classified profile-aware E2E architecture as CROSS since the profile system is shared across features
- All command renaming, Cobra grouping, and migration-specific specs classified as LOCAL since they are specific to the forge-cli-v3 migration
- No overlap detection needed -- docs/decisions/ is empty and no docs/lessons/ directory exists

## Test Results
- **Tests Executed**: No (noTest task)
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] Business rules extracted from PRD into biz-specs.md preview
- [x] Technical specs extracted from PRD into tech-specs.md preview
- [x] CROSS items integrated to project-level dirs
- [x] Integration marker written
- [x] No overlap with existing project specs (verified empty dirs)

## Notes
No tech-design.md exists for this feature; all extractions sourced from prd-spec.md. docs/decisions/ directory exists but is empty. No docs/lessons/ directory. No docs/business-rules/ or docs/conventions/ existed prior to this task.
