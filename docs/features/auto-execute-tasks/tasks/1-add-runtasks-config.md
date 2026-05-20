---
id: "1"
title: "Add RunTasks ModeToggle to AutoConfig"
priority: "P0"
estimated_time: "20m"
dependencies: []
scope: "backend"
breaking: false
type: "coding.feature"
mainSession: false
---

# 1: Add RunTasks ModeToggle to AutoConfig

## Description

在 Go CLI 的 `AutoConfig` 结构体中添加 `RunTasks ModeToggle` 字段，控制 `/quick` 流水线是否跳过 Step 2 确认门直接进入任务生成+执行。

默认值：`quick: true, full: false`——quick 模式默认自动执行，full 模式保留确认门。

## Reference Files
- `docs/proposals/auto-execute-tasks/proposal.md` — Source proposal
- `forge-cli/pkg/forgeconfig/config.go` — AutoConfig 结构体定义与默认值
- `forge-cli/pkg/forgeconfig/config_test.go` — AutoConfig 测试

## Acceptance Criteria

- [ ] `AutoConfig` 结构体新增 `RunTasks ModeToggle` 字段（`yaml:"runTasks"`）
- [ ] `AutoConfigDefaults()` 中设置 `RunTasks: ModeToggle{Quick: true, Full: false}`
- [ ] `IsZero()` 检查包含 `RunTasks` 字段
- [ ] `WithDefaults()` 处理 `RunTasks` 的零值填充
- [ ] `applyDefaults()` 处理 `RunTasks` 的 raw/默认值逻辑
- [ ] `forge config get auto.runTasks` 正确返回 `quick: true, full: false`（默认）
- [ ] 向后兼容：未配置 `runTasks` 时使用默认值
- [ ] 现有测试通过，新增字段有对应测试覆盖

## Hard Rules

- 复用已有的 `ModeToggle` 模式，不引入新的配置结构
- 默认值必须为 `Quick: true, Full: false`

## Implementation Notes

- 参照 `CleanCode` 或 `Validation` 字段的实现模式
- `applyDefaults` 中需要处理 `runTasks` 的 per-mode default
- `config.go` 中的 `getAutoValue` 函数需要添加 `runTasks` case（如有）
