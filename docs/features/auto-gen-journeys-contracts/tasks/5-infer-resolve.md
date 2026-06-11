---
id: "5"
title: "更新 infer.go 和 ResolveFirstTestDep 适配新任务类型"
priority: "P1"
estimated_time: "1.5h"
dependencies: ["3", "4"]
scope: "backend"
breaking: false
type: "coding.enhancement"
mainSession: false
---

# 5: 更新 infer.go 和 ResolveFirstTestDep 适配新任务类型

## Description

更新 `InferType()` 和 `ResolveFirstTestDep()` 函数，使其能正确识别和处理新增的 gen-journeys 和 gen-contracts 任务类型。

## Reference Files
- `docs/proposals/auto-gen-journeys-contracts/proposal.md` — Source proposal
- `forge-cli/pkg/task/infer.go` — InferType (L10-43)
- `forge-cli/pkg/task/autogen.go` — ResolveFirstTestDep (L529-582)

## Acceptance Criteria

- [ ] `InferType()` 识别 `T-test-gen-journeys`（含 type suffix 变体如 T-test-gen-journeys-tui）返回正确的类型
- [ ] `InferType()` 识别 `T-test-gen-contracts` 返回正确的类型
- [ ] `ResolveFirstTestDep()` Breakdown 分支：首任务从 T-eval-journey 更新为 T-test-gen-journeys（通过 findTaskIndexByPrefix 查找）
- [ ] `ResolveFirstTestDep()` Quick 分支：首任务从 T-quick-gen-and-run 更新为 T-test-gen-journeys（通过 findTaskIndexByPrefix 查找）
- [ ] `ExtractTypeSuffix()` 正确处理 T-test-gen-journeys-{type} 的 type suffix 提取
- [ ] 所有现有单测通过

## Hard Rules

- InferType 的模式匹配顺序：更具体的模式（gen-journeys）必须在更通用的模式（gen-scripts）之前
- ResolveFirstTestDep 中 findTaskIndexByPrefix 返回 -1 时 panic 并输出明确错误信息

## Implementation Notes

- InferType 目前使用 if-else 链匹配 ID 前缀。新增的 gen-journeys/gen-contracts 模式需插入到正确的位置
- ResolveFirstTestDep 的 Quick 分支当前查找 `T-quick-gen-and-run`（L567），需更新为查找 `T-test-gen-journeys`
- 注意保持向后兼容：InferType 仍需识别旧的 T-quick-gen-and-run 和 T-test-gen-scripts 模式
