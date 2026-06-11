---
created: "2026-05-27"
tags: [architecture]
---

# Review 任务依赖缺失导致在被审查任务之前执行

## Problem

T-review-doc（P1，依赖 [9]）在 task 11（P2，依赖 [10]）之前执行。T-review-doc 应该审查 task 11 的产出（OVERVIEW/WORKFLOW 文档），但 task 11 还没执行就被审查了。审查结果中 task 11 的 AC 全部标记为"scope 外"。

## Root Cause

因果链（3 层）：

1. **表面现象**：task 10 完成后，T-review-doc（P1）和 task 11（P2）同时 unblocked，claim 按优先级排序选了 T-review-doc
2. **直接原因**：T-review-doc 的 dependencies 只有 [9]，不包含 [10]、[11]。而 task 11 是 doc 类型任务，属于 T-review-doc 的审查范围
3. **根因**：`forge task index` 自动生成 T-review-doc 时，依赖计算只基于"哪些 doc 任务已完成"（task 9 是最后一个手动创建的 doc 任���），没有考虑"还有哪些 doc 任务尚未创建但将被创建"。review-doc 应该依赖所有 doc 类型任务（7, 8, 9, 10, 11, 12），但生成时 task 10-12 可能还未被加入依赖列表

## Solution

`forge task index` 生成 review-doc 时，依赖应包含所有 doc 类型任务的 ID：

```
dependencies: [所有 type.startsWith("doc") 的任务 ID]
```

而非硬编码为最后一个已知的 doc 任务 ID。

## Reusable Pattern

- **Review 任务必须依赖所有被审查任务**：review-doc / review-code 等 auto-generated review 任务的依赖列表应该是"所有同类型任务"的 ID 列表，而非某个固定任务 ID
- **Priority 排序会放大依赖缺失的影响**：当 review-doc 是 P1 而被审查任务是 P2 时，依赖缺失必然导致 review 先执行。如果优先级相同，执行顺序不确定
- **Auto-generated 任务的依赖需要前瞻**：生成 review-doc 时不能只看"当前已有哪些任务"，要看"最终会有哪些 doc 任务"。最安全的做法是让 review-doc 依赖所有 doc 任务

## Related Files

- `forge-cli/pkg/task/autogen.go` — review-doc 依赖生成逻辑
- [[gotcha-task-reference-files-scope-creep]] — task scope 边界问题
