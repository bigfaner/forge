---
status: "completed"
started: "2026-05-26 03:58"
completed: "2026-05-26 04:05"
time_spent: "~7m"
---

# Task Record: T-specs-consolidate Consolidate Specs

## Summary
Extracted 17 rules/specs from surface-aware-justfile feature docs, auto-integrated 11 CROSS items into 3 new project-level spec files, fixed 2 drift issues in existing specs, regenerated vocabulary index

## Changes

### Files Created
- docs/business-rules/surface-orchestration.md
- docs/conventions/surface-cli.md
- docs/conventions/surface-rules.md
- docs/features/surface-aware-justfile/specs/biz-specs.md
- docs/features/surface-aware-justfile/specs/tech-specs.md
- docs/features/surface-aware-justfile/specs/review-choices.md
- docs/features/surface-aware-justfile/specs/.integrated

### Files Modified
- docs/business-rules/error-reporting.md
- docs/conventions/forge-cli-reference.md
- docs/features/surface-aware-justfile/manifest.md
- docs/.vocabulary.md

### Key Decisions
无

## Document Metrics
17 rules extracted (9 biz + 8 tech), 11 CROSS auto-integrated to 3 new spec files, 2 drift fixes applied (error-reporting.md added MIGRATION_REQUIRED, forge-cli-reference.md added --json/--types flags), vocabulary index regenerated

## Referenced Documents
- docs/features/surface-aware-justfile/prd/prd-spec.md
- docs/features/surface-aware-justfile/prd/prd-user-stories.md
- docs/features/surface-aware-justfile/design/tech-design.md
- docs/features/surface-aware-justfile/manifest.md
- docs/business-rules/error-reporting.md
- docs/business-rules/quality-gate.md
- docs/business-rules/task-lifecycle.md
- docs/conventions/error-handling.md
- docs/conventions/forge-cli-reference.md
- docs/conventions/forge-distribution.md

## Review Status
completed

## Acceptance Criteria
- [x] All feature documents scanned for extractable rules
- [x] Preview files generated in specs/ directory
- [x] CROSS items auto-integrated to project-level dirs (non-interactive mode)
- [x] Integration marker written (.integrated)
- [x] Existing project-level specs checked for drift
- [x] Drift fixes applied (error-reporting.md, forge-cli-reference.md)
- [x] Vocabulary index regenerated
- [x] Changes committed with [auto-specs] tag

## Notes
Non-interactive mode (run-tasks dispatcher). All CROSS items auto-integrated. No overlaps with existing decisions or lessons detected. Two drift fixes: BIZ-error-reporting-001 updated to include MIGRATION_REQUIRED exit code 2 case, forge-cli-reference.md updated surfaces command entry with --json and --types flags.
