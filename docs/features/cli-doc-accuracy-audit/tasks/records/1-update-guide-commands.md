---
status: "completed"
started: "2026-06-07 21:45"
completed: "2026-06-07 21:47"
time_spent: "~2m"
---

# Task Record: 1 Update existing command descriptions in guide.md

## Summary
Fixed 4 CLI command reference errors in guide.md: G1 (validate-index -> validate [file]), G2 (quality-gate description), G3 (cleanup description), G4 (task submit --quiet flag)

## Changes

### Files Created
无

### Files Modified
- plugins/forge/hooks/guide.md

### Key Decisions
无

## Document Metrics
4 command descriptions corrected, verified against RunE source code

## Referenced Documents
- docs/proposals/cli-doc-accuracy-audit/proposal.md
- docs/conventions/forge-distribution.md

## Review Status
final

## Acceptance Criteria
- [x] guide.md 中 forge task validate-index 替换为 forge task validate [file]
- [x] guide.md 中 forge quality-gate 描述准确反映实际行为（含 fix task 自动创建、retry-once、docs-only 跳过）
- [x] guide.md 中 forge cleanup 描述从 'clean stale artifacts' 改为具体行为说明（包含 blocked/suspended/rejected 状态的清理）
- [x] guide.md 中 forge task submit 描述补充 --quiet 标志

## Notes
All descriptions verified against Go source code: validate.go (Use: 'validate [file]'), quality_gate.go (AddFixTask + runUnitTestStep retry-once + IsDocsOnly), cleanup.go (StatusCompleted/Blocked/Suspended/Rejected), submit.go (--quiet flag)
