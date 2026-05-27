---
id: "6"
title: "Surface 硬性约束与 forge init 集成"
priority: "P1"
estimated_time: "1h"
complexity: "medium"
dependencies: [4, 5]
surface-key: "."
surface-type: "cli"
breaking: true
type: "coding.feature"
mainSession: false
---

# 6: Surface 硬性约束与 forge init 集成

## Description
实现 Surface 推断失败时的硬性报错机制。`quality_gate.go` 的 `addSingleFixTask()` 中 `inferSurface()` 失败时返回错误——这是行为变更（当前为静默空值）。`forge task add` 命令路径保留软性行为。在 `forge init` 中集成 surface 配置步骤（`forge surfaces detect`），确保新项目初始化即具备 surfaces 配置。硬性失败逻辑独立封装为 `requireSurfaceInference()` 函数。

## Reference Files
- `forge-cli/pkg/task/quality_gate.go`: addSingleFixTask() 中 inferSurface() 失败时返回错误 (source: proposal.md#Surface-硬性约束)
- `forge-cli/cmd/init.go` 或对应 init 命令文件: 集成 forge surfaces detect 步骤 (source: proposal.md#Surface-硬性约束)
- `forge-cli/pkg/task/add.go`: forge task add 路径保持软性行为（推断失败用空字符串） (source: proposal.md#Surface-硬性约束)

## Acceptance Criteria
- [ ] `quality_gate.go` 的 `addSingleFixTask()` 中 surface 推断失败时返回错误，不再创建空 surface 的任务文件
- [ ] `forge task add` 命令路径保持软性行为——推断失败时使用空字符串而非报错
- [ ] `forge init` 中包含 surface 配置步骤（`forge surfaces detect`），新项目初始化时具备 surfaces 配置
- [ ] 硬性失败逻辑封装为独立的 `requireSurfaceInference()` 函数，错误信息包含 `forge surfaces detect` 命令指引

## Implementation Notes
- 硬性失败是行为变更（当前行为为静默空值），而非纯重构
- 硬性失败与模板引擎迁移解耦——即使回滚硬性失败，模板引擎迁移不受影响
- 回滚策略：`requireSurfaceInference()` 可通过单行 revert 恢复为软性行为

### Test Impact
- Affected test suite(s): `forge-cli/pkg/task/`, `forge-cli/cmd/`
- Expected fixture changes: quality gate test fixtures (addSingleFixTask error cases)
- Risk level: medium
