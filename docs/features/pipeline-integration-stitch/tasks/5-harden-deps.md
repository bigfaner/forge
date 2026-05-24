---
id: "5"
title: "加固 mixed feature 依赖注入 + 更新 findFirstTestTaskIdx"
priority: "P1"
estimated_time: "1h"
dependencies: ["1"]
scope: "backend"
breaking: false
type: "coding.enhancement"
mainSession: false
---

# 5: 依赖注入加固

## Description

加固 `build.go` 中 mixed feature 依赖注入逻辑。当前 `ResolveFirstTestDep` 设置依赖后才 prepend `T-review-doc`，存在顺序耦合。将 T-review-doc 依赖和 ResolveFirstTestDep 合并为单次处理。更新 `findFirstTestTaskIdx` 中 quick-mode fallback：移除 `T-quick-gen-and-run*` 匹配，替换为新 test task 前缀（`T-test-gen-journeys` 或 `T-test-gen-contracts`）。

## Reference Files
- `docs/proposals/pipeline-integration-stitch/proposal.md` — Source proposal
- `forge-cli/pkg/task/build.go` — 依赖注入和 findFirstTestTaskIdx 逻辑

## Acceptance Criteria
- [ ] T-review-doc 依赖和 ResolveFirstTestDep 合并为单次处理，无顺序耦合
- [ ] mixed feature re-index 幂等：多次 BuildIndex 后 T-review-doc 依赖不丢失
- [ ] `findFirstTestTaskIdx` quick-mode fallback 使用新 test task 前缀匹配（`T-test-gen-journeys*` 或 `T-test-gen-contracts`）
- [ ] `findFirstTestTaskIdx` 对 quick-mode pipeline 正确返回首个 test task 索引（非 -1）
- [ ] 所有现有测试通过

## Hard Rules
- 依赖图必须保持 DAG（无环）
- findFirstTestTaskIdx 必须同时支持 breakdown 和 quick 两种模式

## Implementation Notes
- 当前 `findFirstTestTaskIdx`（build.go ~line 485）quick-mode fallback 检查 `T-quick-gen-and-run*`，移除后需替换为新的匹配模式
