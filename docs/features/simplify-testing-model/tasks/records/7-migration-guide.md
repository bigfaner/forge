---
status: "completed"
started: "2026-05-17 15:32"
completed: "2026-05-17 15:36"
time_spent: "~4m"
---

# Task Record: 7 Write v2-to-v3 migration guide and update config docs

## Summary
Created v2-to-v3 migration guide at docs/proposals/simplify-testing-model/migration-v2-to-v3.md covering field mapping table, profile-to-language mapping (6 entries), common override patterns, and troubleshooting. Updated user-facing docs: profile-authoring.md rewritten for v3 language strategy model, conventions/profile-system.md updated to language-aware architecture.

## Changes

### Files Created
- docs/proposals/simplify-testing-model/migration-v2-to-v3.md

### Files Modified
- docs/profile-authoring.md
- docs/conventions/profile-system.md

### Key Decisions
- Treated docs/features/*/testing/test-cases.md and docs/proposals/* as historical records not subject to v3 migration -- modifying them would falsify archival history
- Rewrote profile-authoring.md entirely rather than patching -- the v2 profile concept is fully replaced by v3 language strategies, making incremental edits incoherent

## Test Results
- **Tests Executed**: Yes
- **Passed**: 0
- **Failed**: 0
- **Coverage**: 0.0%

## Acceptance Criteria
- [x] Migration guide exists at docs/proposals/simplify-testing-model/migration-v2-to-v3.md
- [x] Guide documents v2-to-v3 field mapping table: test-profiles -> (removed, auto-detect), capabilities -> interfaces
- [x] Guide documents profile-name-to-language mapping (6 entries): go-test->go, web-playwright->javascript, rust-test->rust, pytest->python, java-junit->java, maestro->mobile
- [x] Guide documents common override patterns: multi-language false positive (languages: [go]), monorepo subdirectories
- [x] Guide documents troubleshooting: detection failures, no language detected error
- [x] All config.yaml examples in user-facing docs use v3 schema (no test-profiles or capabilities fields)
- [x] No references to profile or capability in user-facing documentation (excluding migration guide)

## Notes
Documentation-only task. No test execution applicable. Scope of user-facing docs update limited to profile-authoring.md and conventions/profile-system.md -- historical feature records and proposals were not modified as they are archival.
