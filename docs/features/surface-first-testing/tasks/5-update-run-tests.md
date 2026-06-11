---
id: "5"
title: "更新 run-tests Convention 读取路径"
priority: "P1"
estimated_time: "1h"
dependencies: [2]
type: "doc"
mainSession: false
---

# 5: 更新 run-tests Convention 读取路径

## Description
将 run-tests 的 Convention 读取路径从旧路径改为 `testing/{surface}/`，读取 per-surface 的 core.md 获取超时策略和生命周期规则。增加旧结构检测和迁移提示输出。

## Reference Files
- `docs/proposals/surface-first-testing/proposal.md` — Proposed Solution (下游 skill 适配), Non-Functional Requirements
- `plugins/forge/skills/run-tests/SKILL.md`: Convention 读取路径需修改 (ref: Proposed Solution)

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/run-tests/SKILL.md` | Convention 读取路径改为 `testing/{surface}/` |

## Acceptance Criteria
- [ ] SKILL.md Convention 读取路径改为 `testing/{surface}/core.md`
- [ ] 检测到旧结构 convention 文件时输出迁移提示而非静默失败

## Hard Rules
- 必须先加载 `docs/conventions/forge-distribution.md` 了解 plugin 分发模型

## Implementation Notes
- run-tests 从 core.md 读取的关键信息：超时策略、生命周期规则、隔离模型
- 旧结构检测逻辑与 gen-test-scripts 一致：检查 `docs/conventions/testing/` 下是否存在旧格式文件
