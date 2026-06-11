---
status: "completed"
started: "2026-06-07 22:06"
completed: "2026-06-07 22:08"
time_spent: "~2m"
---

# Task Record: T-review-doc Review Documentation Quality

## Summary
Reviewed guide.md against 8 acceptance criteria from Task 1 (update-guide-commands) and Task 2 (add-guide-entries). All CLI commands and flags verified against forge --help output and source code. All 8 ACs PASS — no fixes required.

## Changes

### Files Created
无

### Files Modified
无

### Key Decisions
无

## Document Metrics
8 ACs verified: 8 PASS, 0 FAIL. Cross-referenced against CLI --help output and Go source code (cleanup.go).

## Referenced Documents
- docs/proposals/cli-doc-accuracy-audit/proposal.md

## Review Status
reviewed

## Acceptance Criteria
- [x] guide.md 中 forge task validate-index 替换为 forge task validate [file]
- [x] guide.md 中 forge quality-gate 描述含 fix task 自动创建、retry-once、docs-only 跳过
- [x] guide.md 中 forge cleanup 描述含 blocked/suspended/rejected 状态清理
- [x] guide.md 中 forge task submit 描述补充 --quiet 标志
- [x] guide.md 新增 forge task query 命令描述与 --help 一致
- [x] guide.md 新增 forge task check-deps 命令描述与 --help 一致
- [x] guide.md 新增 forge feature list 命令描述与 --help 一致
- [x] guide.md 中 forge task list 描述补充 --tree 标志

## Notes
No modifications needed — all 8 ACs already satisfied by current guide.md content. Verified forge cleanup behavior against Go source (cleanup.go line 70) confirming it handles completed/blocked/suspended/rejected statuses as documented.
