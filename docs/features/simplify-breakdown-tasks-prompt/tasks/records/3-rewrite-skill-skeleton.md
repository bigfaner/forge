---
status: "completed"
started: "2026-05-19 23:05"
completed: "2026-05-19 23:11"
time_spent: "~6m"
---

# Task Record: 3 Rewrite SKILL.md skeleton with condition-rule matrix

## Summary
Rewrote SKILL.md from 421 lines/23KB to 144 lines/7.7KB. Replaced 6 conditional tags (HAS_UI, NO_UI, UI_ONLY, HAS_PLACEMENT, RULE, HAS_DB) with a 4-row Condition-Rule Matrix as the first instruction block after Prerequisites. Each step now uses IF-loaded checks to conditionally apply rule files. Skeleton is structurally complete without any rule files — fail-safe design preserved.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/breakdown-tasks/SKILL.md

### Key Decisions
- Compressed Scope Assignment algorithm into inline paragraph format to meet size target while preserving all classification logic
- Removed mermaid diagram from Docs-Only Fast Path (not in acceptance criteria as required)
- Kept ui-placement OR condition per explicit acceptance criteria spec despite Hard Rules about no boolean expressions (AC overrides)

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] SKILL.md reduced from ~23KB to <=8KB
- [x] All 6 conditional tags removed
- [x] Condition-rule matrix added as first instruction block after Prerequisites
- [x] Matrix contains exactly 4 rows with specified conditions
- [x] Each step re-prints load instruction inline as safety net
- [x] Skeleton never contains rule content inline
- [x] Always-needed rules preserved inline
- [x] Prerequisites section updated (tags removed, artifact table kept)
- [x] Step 1 updated with IF rule file loaded instructions
- [x] Step 2 updated with base mapping inline, conditional via rule files
- [x] Step 3 updated with IF loaded apply, else artifact-driven
- [x] Step 4a updated with scope/type/template always inline
- [x] Steps 5-7 unchanged
- [x] All paths use skill-relative references via CLAUDE_SKILL_DIR
- [x] Skeleton structurally complete without any rule files

## Notes
Final size: 144 lines / 7659 bytes. All 4 rule files (created in tasks 1-2) correctly referenced. ui-placement condition uses OR per acceptance criteria spec line 50.
