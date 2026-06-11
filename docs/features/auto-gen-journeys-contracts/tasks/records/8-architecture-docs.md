---
status: "completed"
started: "2026-05-24 00:44"
completed: "2026-05-24 00:45"
time_spent: "~1m"
---

# Task Record: 8 更新 ARCHITECTURE.md 和相关文档

## Summary
Updated ARCHITECTURE.md test pipeline section with Breakdown/Quick mode chains and staged-across-types topology; added deprecated comment to test-gen-and-run.md

## Changes

### Files Created
无

### Files Modified
- docs/ARCHITECTURE.md
- forge-cli/pkg/task/data/test-gen-and-run.md

### Key Decisions
无

## Document Metrics
2 files modified: ARCHITECTURE.md (test pipeline section rewritten ~15→68 lines), test-gen-and-run.md (1 line added)

## Referenced Documents
- docs/proposals/auto-gen-journeys-contracts/proposal.md

## Review Status
final

## Acceptance Criteria
- [x] ARCHITECTURE.md Breakdown mode chain: gen-journeys → eval-journey → gen-contracts → eval-contract → gen-scripts → run → verify
- [x] ARCHITECTURE.md Quick mode chain: gen-journeys → gen-contracts → gen-scripts → run → verify (no eval gates)
- [x] test-gen-and-run.md header has deprecated comment
- [x] Documentation describes staged across types topology

## Notes
Hard Rules respected: test-gen-and-run.md not deleted, only deprecated comment added; template content unchanged
