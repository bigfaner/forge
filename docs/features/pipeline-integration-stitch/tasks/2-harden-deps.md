---
id: "2"
title: "加固依赖注入：合并 resolveTestDepsAndInjectReviewDoc + 更新 findFirstTestTaskIdx"
priority: "P1"
estimated_time: "1.5h"
dependencies: ["1"]
scope: "backend"
breaking: true
type: "coding.enhancement"
mainSession: false
---

# 2: 加固依赖注入 + 更新 findFirstTestTaskIdx

## Description

修复 P1 中两个 build.go 问题：(1) T-review-doc prepend 与 ResolveFirstTestDep 顺序耦合——两步操作硬编码顺序，重排会导致 T-review-doc 丢失；(2) findFirstTestTaskIdx quick-mode 分支匹配废弃类型 `T-quick-gen-and-run*`，当前靠 fallback 意外正确。

## Reference Files
- `proposal.md#P1-—-CategoryEval-+-依赖加固` — 依赖注入合并和 findFirstTestTaskIdx 更新规格
- `proposal.md#Key-Risks` — findFirstTestTaskIdx 修改对 dependency wiring 的影响风险
- `proposal.md#Success-Criteria` — resolveTestDepsAndInjectReviewDoc 的行为验证标准

## Acceptance Criteria

- [ ] `resolveTestDepsAndInjectReviewDoc(testTasks, idx, "quick", true)` 返回的依赖列表包含 T-review-doc
- [ ] `resolveTestDepsAndInjectReviewDoc(testTasks, idx, "quick", false)` 返回的依赖列表不包含 T-review-doc 且与旧 ResolveFirstTestDep 输出一致
- [ ] `findFirstTestTaskIdx` quick-mode 分支使用 `findTaskIndexByPrefix(tasks, "T-test-gen-journeys")`
- [ ] BuildIndex 中不再有独立的 T-review-doc prepend 操作（已合并到新函数）
- [ ] 集成测试覆盖 Quick mode 完整依赖链：创建 Quick mode feature → needsEval=true 时任务列表包含 T-review-doc → gen-journeys 依赖指向 T-review-doc
- [ ] 所有现有测试通过

## Hard Rules

- 合并函数签名：`func resolveTestDepsAndInjectReviewDoc(testTasks []AutoGenTaskDef, index *TaskIndex, mode string, needsEval bool)`
- findFirstTestTaskIdx 和 ResolveFirstTestDep 使用相同发现机制：`findTaskIndexByPrefix(tasks, "T-test-gen-journeys")`
- T-review-doc 在依赖链中必须排在首位（它依赖 last business task，original dep 通过 T-review-doc 传递）

## Implementation Notes

**findFirstTestTaskIdx 的 P1/P2 重叠**：此任务中更新 findFirstTestTaskIdx 的 quick-mode 匹配逻辑（P1），与 Task 3 中清理同一位置的 gen-and-run 注释和残留代码（P2）是同一处代码。此任务完成 P1 部分的修改，Task 3 负责 P2 部分的注释清理。

**集成测试位置**：扩展现有 `pkg/task/build_test.go` 中的 `TestBuildIndex` 系列测试。
