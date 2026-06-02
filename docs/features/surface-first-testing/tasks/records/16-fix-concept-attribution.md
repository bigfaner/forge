---
status: "completed"
started: "2026-06-02 23:45"
completed: "2026-06-02 23:47"
time_spent: "~2m"
---

# Task Record: 16 Fix gen-contracts concept attribution + run-tests cross-skill reference

## Summary
Fixed gen-contracts SKILL.md concept attribution (3.7 = TUI-specific, 3.8-3.10 = all surfaces) and removed run-tests env-check.md cross-skill reference to gen-journeys, replaced with self-contained reference to run-tests' own rules/surfaces/<type>.md

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/gen-contracts/SKILL.md
- plugins/forge/skills/run-tests/rules/env-check.md

### Key Decisions
无

## Document Metrics
2 files modified, 4 acceptance criteria met, 0 cross-skill references remaining

## Referenced Documents
- docs/proposals/surface-first-testing/proposal.md
- docs/conventions/forge-distribution.md
- plugins/forge/skills/gen-contracts/rules/tui-async.md

## Review Status
final

## Acceptance Criteria
- [x] gen-contracts SKILL.md Step 3.7 clearly annotated as TUI-specific
- [x] gen-contracts SKILL.md Steps 3.8/3.9/3.10 clearly annotated as applies to all surface types
- [x] run-tests env-check.md no longer references gen-journeys skill internal file paths
- [x] env-check.md references run-tests own rules/surfaces/<type>.md or is self-contained

## Notes
env-check.md now fully self-contained with per-surface detection items inline and references to run-tests' own rules/surfaces/ directory. gen-contracts SKILL.md Steps 3.7-3.10 now have clear scope annotations.
