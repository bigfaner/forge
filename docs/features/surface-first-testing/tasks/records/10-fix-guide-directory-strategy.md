---
status: "completed"
started: "2026-06-02 22:49"
completed: "2026-06-02 22:50"
time_spent: "~1m"
---

# Task Record: 10 Fix guide.md stale reference + unify test directory strategy

## Summary
Fixed guide.md stale test-type-model.md link (removed, made self-contained). Unified test directory strategy across guide.md, web.md, mobile.md from tests/e2e/ to tests/<journey>/, consistent with gen-test-scripts HARD-RULE.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/hooks/guide.md
- plugins/forge/skills/test-guide/templates/surfaces/web.md
- plugins/forge/skills/test-guide/templates/surfaces/mobile.md

### Key Decisions
无

## Document Metrics
Testing section: 12 lines (within 20-line budget). 3 files modified, 0 stale references remaining.

## Referenced Documents
- docs/proposals/surface-first-testing/proposal.md
- plugins/forge/skills/gen-test-scripts/SKILL.md

## Review Status
final

## Acceptance Criteria
- [x] guide.md no longer contains docs/reference/test-type-model.md link, mapping table and e2e constraint are self-contained
- [x] guide.md web/mobile test directory changed from tests/e2e/ to tests/<journey>/, consistent with gen-test-scripts
- [x] test-guide/templates/surfaces/web.md file location no longer references tests/e2e/
- [x] test-guide/templates/surfaces/mobile.md file location no longer references tests/e2e/
- [x] guide.md Testing section total lines <= 20

## Notes
Testing section is 12 lines (was 13 before, net reduction of 1 line). Hard Rules satisfied: guide.md delta <= 20 lines, no references to other skill internal files.
