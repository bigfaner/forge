---
status: "completed"
started: "2026-06-07 00:15"
completed: "2026-06-07 00:31"
time_spent: "~16m"
---

# Task Record: fix-2 Fix: justfile-integration 7 pre-existing failures

## Summary
Fixed all justfile-integration test failures by updating assertions to match evolved SKILL.md structure, project-detection.md delegation, mixed.just per-surface recipe model, and renamed custom recipes (claude-p replaces claude/claude-c). Also fixed a hidden panic in TC_DET_019 that was masking 8 additional test failures beyond the original 7.

## Changes

### Files Created
无

### Files Modified
- tests/justfile-integration/forge_detection_test.go
- tests/justfile-integration/init_justfile_test.go
- tests/justfile-integration/mixed_cli_test.go

### Key Decisions
- Detection signal tests (TC_DET_001-004, 008, 019) now check rules/project-detection.md instead of SKILL.md directly, matching the delegation pattern in Step 1a
- Mixed template tests (TC_MIX_002-012) updated from scope=""/case-esac model to per-surface frontend-*/backend-* recipe model
- Custom recipe test (TC_FJ_010, TC_020) accepts any custom recipe outside boundary markers instead of hardcoding claude:/claude-c:
- CLI integration tests (TC_001, TC_002, TC_007, TC_008, TC_014) updated to match current agent/command content (submit-task replaces record-task, no Breaking Task Gate section)

## Test Results
- **Tests Executed**: Yes
- **Passed**: 90
- **Failed**: 0
- **Coverage**: 0.0%

## Acceptance Criteria
- [x] All justfile-integration tests pass (0 failures)
- [x] Test assertions match current SKILL.md and justfile content
- [x] No panic/crash in test suite

## Notes
The original 7 failures were caused by: (1) SKILL.md detection logic moved to rules/project-detection.md, (2) custom recipes renamed from claude:/claude-c: to claude-p:/install-forge:, (3) SKILL.md section renamed from '## Workflow' to '## Process Flow'. Fixing the TC_DET_019 panic revealed 8 additional hidden failures caused by mixed.just evolving from scope dispatch (case/esac) to per-surface recipes (frontend-*/backend-*).
