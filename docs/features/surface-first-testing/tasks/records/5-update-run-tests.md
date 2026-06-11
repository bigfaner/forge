---
status: "completed"
started: "2026-06-02 22:02"
completed: "2026-06-02 22:03"
time_spent: "~1m"
---

# Task Record: 5 更新 run-tests Convention 读取路径

## Summary
Updated run-tests SKILL.md: added Step 2.5 (Load Convention) to read per-surface timeout and lifecycle rules from testing/{surface}/core.md, with legacy structure detection and migration prompt output

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/run-tests/SKILL.md

### Key Decisions
无

## Document Metrics
1 file modified, 4 sections changed (workflow diagram, new Step 2.5, Step 5 reference, Error Handling table)

## Referenced Documents
- docs/proposals/surface-first-testing/proposal.md
- docs/conventions/forge-distribution.md
- plugins/forge/skills/gen-test-scripts/SKILL.md

## Review Status
final

## Acceptance Criteria
- [x] SKILL.md Convention read path changed to testing/{surface}/core.md
- [x] Legacy structure detection outputs migration prompt instead of silent failure

## Notes
Legacy detection logic aligned with gen-test-scripts SKILL.md pattern (same wording, same exit code 2). Convention loading is graceful: missing core.md falls through to orchestration rule defaults.
