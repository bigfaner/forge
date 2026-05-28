---
id: "9"
title: "Skill/Command/Agent 概念对齐与 scope 清理"
priority: "P2"
estimated_time: "1h"
complexity: "low"
dependencies: [8]
surface-key: "."
surface-type: "cli"
breaking: false
type: "coding.refactor"
mainSession: false
---

# 9: Skill/Command/Agent 概念对齐与 scope 清理

## Description
清除 Forge 插件中所有 deprecated `scope` 概念残留。删除已完成的迁移文档 `scope-to-surface-key.md`。清理 `mixed.just` 模板中的 `frontend`/`backend` scope 参数，改用 surface-aware recipe 命名。审查 `task-executor.md` 和 `commands/*.md` 确认无 deprecated scope 字段引用（prose 用法如 "well-scoped" 不算）。`FrontmatterData` struct 的 `scope` 字段标记 `// Deprecated:` 注释。

## Reference Files
- `plugins/forge/skills/breakdown-tasks/rules/scope-to-surface-key.md`: 删除此文件（迁移已完成） (source: proposal.md#Skill/Command/Agent-概念对齐)
- `plugins/forge/skills/init-justfile/templates/mixed.just`: 清除 frontend/backend scope 参数 (source: proposal.md#Skill/Command/Agent-概念对齐)
- `plugins/forge/agents/task-executor.md`: 审查无 scope 残留 (source: proposal.md#Skill/Command/Agent-概念对齐)
- `plugins/forge/commands/*.md`: 审查无 scope 残留 (source: proposal.md#Skill/Command/Agent-概念对齐)

## Acceptance Criteria
- [ ] `scope-to-surface-key.md` 已删除
- [ ] `mixed.just` 模板中无 `scope` 参数和 `frontend`/`backend` 值，改用 surface-aware recipe 命名
- [ ] `task-executor.md` 和 `commands/*.md` 中无 deprecated scope 字段引用（prose 用法如 "well-scoped" 不算）
- [ ] `FrontmatterData` struct 的 `scope` 字段包含 `// Deprecated:` 注释，`CheckLegacyScope()` 保留用于迁移检测

## Implementation Notes
- 复杂度判定覆盖：AC=4 且无 Hard Rules，4 个 ref > 1。但任务仅涉及 doc/config 编辑和简单 Go 注释添加，认定为 low。
- prose 中的 scope 用法（如 "well-scoped"、"scope of this task"）不算 deprecated 残留，仅清理字段引用和参数传递
- `FrontmatterData.Scope` 标记 deprecated 但不删除，`CheckLegacyScope()` 保留用于向后兼容迁移检测
