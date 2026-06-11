---
status: "completed"
started: "2026-05-27 20:02"
completed: "2026-05-27 20:07"
time_spent: "~5m"
---

# Task Record: 10 修复 ARCHITECTURE.md 与实际代码的对齐

## Summary
修复 ARCHITECTURE.md 与实际代码的多处不一致：Quick 模式流程图移除 gen-contracts/gen-scripts；任务 ID 从 T-test-1~5 更新为描述性 ID；移除 T-test-promote 幽灵条目；profile type → surface type/Convention 术语统一；profile 路由 → Convention 路由；并行执行描述按 autogen.go 实际代码修正

## Changes

### Files Created
无

### Files Modified
- docs/ARCHITECTURE.md

### Key Decisions
无

## Document Metrics
1 file modified, 6 acceptance criteria, ~10 edit points

## Referenced Documents
- docs/proposals/test-pipeline-consistency-audit/proposal.md
- forge-cli/pkg/task/autogen.go

## Review Status
final

## Acceptance Criteria
- [x] Quick 模式流程图不含 gen-contracts/gen-scripts
- [x] 任务 ID 为描述性名称（非 T-test-1~5）
- [x] 不含 T-test-promote 条目
- [x] profile type 已替换为 Convention / surface type
- [x] profile 路由已替换为 Convention 路由
- [x] 并行执行描述与 autogen.go 实际代码一致

## Notes
以 autogen.go GetQuickTestTasks() 和 GetBreakdownTestTasks() 为 ground truth 交叉验证所有流程图描述和任务映射
