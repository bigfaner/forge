---
id: "2"
title: "Update lesson with Impact Declaration pattern"
priority: "P2"
estimated_time: "15m"
dependencies: ["1"]
type: "doc"
mainSession: false
---

# 2: Update lesson with Impact Declaration pattern

## Description

更新 `docs/lessons/gotcha-characterization-test-vs-refactoring.md`，在 Solution 部分补充 Impact Declaration 机制的说明，将纯问题描述升级为包含解决方案的完整 lesson。

## Reference Files
- `docs/proposals/refactor-impact-declaration/proposal.md` — Source proposal
- `docs/lessons/gotcha-characterization-test-vs-refactoring.md` — Lesson to update

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `docs/lessons/gotcha-characterization-test-vs-refactoring.md` | 补充 Solution 中的 Impact Declaration 机制描述 |

## Acceptance Criteria

- [ ] Lesson 的 Solution 部分引用了 Impact Declaration 机制（PRESERVE/EVOLVE 分类）
- [ ] Reusable Pattern 部分更新，说明执行器现在会在重构前主动声明影响范围
- [ ] 保留原有问题描述和根因分析不变

## Implementation Notes

只更新 Solution 和 Reusable Pattern 部分，不修改 Problem 和 Root Cause 部分（它们是历史记录）。
