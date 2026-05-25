---
status: "completed"
started: "2026-05-25 00:35"
completed: "2026-05-25 00:38"
time_spent: "~3m"
---

# Task Record: 9 Dead code cleanup (example files and templates)

## Summary
Deleted sitemap-example.json after inlining its content into schema.md; evaluated all 6 init-justfile .just templates (retained due to test dependencies)

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/gen-sitemap/rules/schema.md
- docs/features/v3-release-audit/tasks/9-dead-code-cleanup.md

### Key Decisions
无

## Document Metrics
1 file deleted (sitemap-example.json), 1 file modified (schema.md inlined example), 6 templates evaluated and retained

## Referenced Documents
- docs/proposals/v3-release-audit/proposal.md

## Review Status
final

## Acceptance Criteria
- [x] grep -r 'sitemap-example' plugins/forge/ returns 0
- [x] init-justfile 6 .just templates evaluated, usage status recorded
- [x] Deleted files confirmed no references via grep -r

## Notes
sitemap-example.json was referenced by schema.md -- inlined the example into schema.md as a ## Full Example section before deleting. All 6 .just templates are actively used by forge-cli/tests/justfile-integration/ tests and cannot be deleted safely.
