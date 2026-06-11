---
status: "completed"
started: "2026-05-23 09:53"
completed: "2026-05-23 10:07"
time_spent: "~14m"
---

# Task Record: 4 Wire BodyContext through BuildIndex with proposal/PRD data extraction

## Summary
Wire BodyContext through BuildIndex with proposal/PRD data extraction. Added extractBodyContext() that reads proposal (quick mode) or PRD (breakdown mode) to populate Scope, SuccessCriteria, and AcceptanceCriteria fields. Updated both GenerateTestTaskMD call sites to use populated BodyContext instead of empty struct.

## Changes

### Files Created
- forge-cli/pkg/task/extract.go
- forge-cli/pkg/task/extract_test.go

### Files Modified
- forge-cli/pkg/task/build.go

### Key Decisions
- Extract functions in separate file (extract.go) to keep build.go focused on BuildIndex orchestration
- extractBulletItems handles both plain '- item' and checkbox '- [ ] item' / '- [x] item' formats
- extractCheckboxItems only collects top-level checkboxes (skips indented sub-items)
- BodyContext extraction happens once after mode detection, capabilities read moved up to avoid duplicate forgeconfig.ReadInterfaces call
- Missing proposal/PRD files produce empty BodyContext fields (graceful degradation per hard rules)

## Test Results
- **Tests Executed**: Yes
- **Passed**: 612
- **Failed**: 0
- **Coverage**: 90.1%

## Acceptance Criteria
- [x] BuildIndex() extracts Scope from proposal/PRD '## Scope > ### In Scope' section
- [x] BuildIndex() extracts SuccessCriteria from proposal/PRD '## Success Criteria' section
- [x] BuildIndex() extracts AcceptanceCriteria from PRD '## Acceptance Criteria' section (breakdown mode only)
- [x] FeatureSlug and Mode are populated from existing BuildIndex data
- [x] Interfaces populated from .forge/config.yaml via forgeconfig.ReadInterfaces
- [x] Both GenerateTestTaskMD() call sites in build.go pass populated BodyContext
- [x] Existing tests pass (backward compatible -- empty BodyContext produces same output as before)

## Notes
Coverage at 90.1% exceeds 80% target. 14 new tests added covering extraction functions and BuildIndex integration.
