---
id: "2"
title: "重写 test-guide skill"
priority: "P0"
estimated_time: "2h"
dependencies: [1]
type: "doc"
mainSession: false
---

# 2: 重写 test-guide skill

## Description
重写 test-guide 的 SKILL.md，将主导流程从"框架检测 → 单文件生成"改为"Surface 检测 → per-surface 模板渲染"。框架检测从主导流程降格为辅助步骤，仅用于填充 core.md 断言偏好表。新流程从 `.forge/config.yaml` 读取 Surface 配置，从 `templates/surfaces/*.md` 渲染 per-surface 的 convention 文件，并生成顶层速查表。

## Reference Files
- `docs/proposals/surface-first-testing/proposal.md` — Proposed Solution (核心变更 #4), Success Criteria, Out of Scope
- `plugins/forge/skills/test-guide/SKILL.md`: 现有 skill 定义，需全量重写 (ref: Proposed Solution)
- `plugins/forge/skills/test-guide/rules/signal-detection.md`: 框架检测逻辑需重构为辅助步骤 (ref: Out of Scope)
- `plugins/forge/skills/test-guide/rules/convention-structure.md`: Convention 文件结构规则需更新为 surface-first (ref: Proposed Solution)

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/test-guide/SKILL.md` | 全量重写：Surface 检测 → 模板渲染 → 速查表生成 |
| `plugins/forge/skills/test-guide/rules/convention-structure.md` | 更新为 surface-first 目录结构 |

## Acceptance Criteria
- [ ] SKILL.md 包含读取 `.forge/config.yaml` surfaces 配置的步骤
- [ ] SKILL.md 包含从 `templates/surfaces/*.md` 生成 per-surface convention 文件（index.md + core.md）的步骤
- [ ] SKILL.md 包含生成顶层 `docs/conventions/testing/index.md` 速查表的步骤
- [ ] 框架检测（signal-detection）重构为辅助步骤，结果仅用于填充 core.md 断言偏好表
- [ ] 旧的"框架检测 → 单文件生成"流程已移除

## Hard Rules
- 必须先加载 `docs/conventions/forge-distribution.md` 了解 plugin 分发模型后再修改 SKILL.md

## Implementation Notes
- signal-detection.md 规则文件保留但降格——不再驱动主流程，仅在生成 core.md 时作为辅助步骤调用
- convention-structure.md 需描述新的 surface-first 目录结构（`testing/{surface}/core.md`）
- draft-generation.md 和 pattern-extraction.md 可能需要微调以适配新流程，但非核心变更
