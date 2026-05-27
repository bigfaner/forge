---
status: "completed"
started: "2026-05-26 22:11"
completed: "2026-05-26 22:14"
time_spent: "~3m"
---

# Task Record: 2 更新核心文档测试术语

## Summary
更新 ARCHITECTURE.md、guide.md、task-lifecycle.md 中的测试术语，将泛化的 'e2e 测试'/'高级测试' 替换为 surface-specific 测试类型名称

## Changes

### Files Created
无

### Files Modified
- docs/ARCHITECTURE.md
- plugins/forge/hooks/guide.md
- docs/business-rules/task-lifecycle.md

### Key Decisions
无

## Document Metrics
3 files modified: ARCHITECTURE.md (6 处术语替换 + 2 处引用新增), guide.md (1 条目新增), task-lifecycle.md (1 保留类型列表更新)

## Referenced Documents
- docs/proposals/surface-test-type-model/proposal.md
- docs/reference/test-type-model.md

## Review Status
final

## Acceptance Criteria
- [x] ARCHITECTURE.md 中不再出现将所有生成测试统称为 'e2e 测试' 或 '高级测试' 的表述
- [x] guide.md Terminology 部分新增 Test Type 条目，包含 Surface → Test Type 映射简要说明
- [x] task-lifecycle.md 保留类型列表已包含新的 test 类型名（通配模式 test.gen-scripts.* / test.run.*）
- [x] 所有术语变更引用 docs/reference/test-type-model.md 作为权威定义来源

## Notes
guide.md Test Type 条目紧接 Surface Type 之后，符合 Implementation Notes 要求。task-lifecycle.md 使用通配模式而非逐一列举 5 种 surface，符合 Hard Rules。
