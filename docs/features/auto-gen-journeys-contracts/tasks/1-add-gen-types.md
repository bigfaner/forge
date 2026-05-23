---
id: "1"
title: "新增 test.gen-journeys 和 test.gen-contracts 类型定义"
priority: "P0"
estimated_time: "1h"
dependencies: []
scope: "backend"
breaking: false
type: "coding.feature"
mainSession: false
---

# 1: 新增 test.gen-journeys 和 test.gen-contracts 类型定义

## Description

在 types.go 和 autogen.go 中注册两个新的自动生成任务类型，为后续的 Breakdown/Quick 模式集成提供类型基础。

## Reference Files
- `docs/proposals/auto-gen-journeys-contracts/proposal.md` — Source proposal
- `forge-cli/pkg/task/types.go` — 类型常量、TaskTypeRegistry、SystemTypes
- `forge-cli/pkg/task/autogen.go` — autogenTypeToFile 映射

## Acceptance Criteria

- [ ] `TypeTestGenJourneys` 常量值为 `"test.gen-journeys"`（types.go）
- [ ] `TypeTestGenContracts` 常量值为 `"test.gen-contracts"`（types.go）
- [ ] 两个新类型在 `TaskTypeRegistry` 中有对应条目，label 清晰描述用途
- [ ] 两个新类型在 `SystemTypes` map 中注册（值为 true）
- [ ] `autogenTypeToFile` 映射中添加：`TypeTestGenJourneys → "data/test-gen-journeys.md"`, `TypeTestContracts → "data/test-gen-contracts.md"`
- [ ] `forge -h` 输出的类型列表包含 test.gen-journeys 和 test.gen-contracts

## Hard Rules

- 类型常量按字母序插入到现有 test.* 常量区域
- SystemTypes 注册确保这两个类型被视为系统生成、不可手动创建

## Implementation Notes

- 参考 `TypeEvalJourney`/`TypeEvalContract` 的注册方式，保持一致
- `autogenTypeToFile` 的模板文件在 Task 2 中创建，此处只需添加映射条目
