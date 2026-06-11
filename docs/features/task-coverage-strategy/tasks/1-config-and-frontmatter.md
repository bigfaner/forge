---
id: "1"
title: "Add coverage config schema, parsing, and frontmatter field"
priority: "P1"
estimated_time: "2h"
dependencies: []
scope: "backend"
breaking: false
type: "coding.feature"
mainSession: false
---

# 1: Add coverage config schema, parsing, and frontmatter field

## Description

为 `.forge/config.yaml` 新增 `coverage` 配置段，支持按任务类型设置覆盖率策略（百分比目标或 maintain 模式）。同步扩展 Go CLI 的 config 解析和 task frontmatter 解析，使覆盖率配置可从全局配置和任务级别两个维度控制。

**背景**：当前覆盖率目标仅在 CLAUDE.md 中作为文字指导，对 agent 行为无实际约束。本任务建立配置基础设施，为后续 prompt 注入提供数据来源。

## Reference Files
- `docs/proposals/task-coverage-strategy/proposal.md` — Source proposal
- `forge-cli/pkg/forgeconfig/config.go` — Config struct 与解析逻辑
- `forge-cli/pkg/forgeconfig/detect.go` — Config 读取辅助函数
- `forge-cli/pkg/task/frontmatter.go` — Task frontmatter 解析
- `forge-cli/pkg/task/build.go` — Task 构建与类型推断
- `forge-cli/internal/cmd/testdata/forge-config.example.yaml` — 示例配置

## Acceptance Criteria

- `Config` struct 新增 `Coverage *CoverageConfig` 字段
- `CoverageConfig` 包含按任务类型的策略配置，支持 `percentage`（数字目标）和 `maintain`（保持模式）两种策略
- 内置默认值：`coding.feature` → 80%, `coding.enhancement` / `coding.fix` → 60%, `coding.refactor` / `coding.cleanup` / `coding.clean` → maintain
- 无配置时 `CoverageConfig` 返回内置默认值，不报错不阻塞
- `GetConfigValue` 支持 `coverage.*` dot-notation 查询
- `FrontmatterData` 新增可选 `Coverage` 字段（`*int`，nil 表示使用全局默认）
- `forge-config.example.yaml` 更新包含 coverage 配置示例
- 现有测试通过

## Hard Rules

- 配置格式使用 map 结构（`map[string]CoverageStrategy`）而非固定字段，确保可扩展
- maintain 策略下不需要 `percentage` 字段，只用 `type: maintain`
- YAML 中未知字段必须被静默忽略（兼容现有行为）

## Implementation Notes

- 参考 `AutoConfig` 的默认值填充模式（`AutoConfigDefaults()`），为 `CoverageConfig` 实现类似的默认值机制
- frontmatter 的 `Coverage` 字段为 `*int`（指针），区分"未设置"和"设置为 0"的情况
- `build.go` 中构建 `Task` struct 时需将 frontmatter 的 coverage 值传递到 task 数据中
