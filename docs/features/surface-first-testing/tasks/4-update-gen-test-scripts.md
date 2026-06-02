---
id: "4"
title: "更新 gen-test-scripts Convention 加载路径"
priority: "P1"
estimated_time: "1h"
dependencies: [2]
type: "doc"
mainSession: false
---

# 4: 更新 gen-test-scripts Convention 加载路径

## Description
将 gen-test-scripts 的 Convention 加载逻辑从旧路径（框架文件）改为新的 surface 目录遍历路径 `testing/{surface}/core.md`。增加旧结构检测和迁移提示输出。确保生成的测试代码使用 per-surface build tag 命名。

## Reference Files
- `docs/proposals/surface-first-testing/proposal.md` — Proposed Solution (下游 skill 适配), Success Criteria, Out of Scope
- `plugins/forge/skills/gen-test-scripts/SKILL.md`: Convention 加载逻辑需修改 (ref: Proposed Solution)

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/gen-test-scripts/SKILL.md` | Convention 加载路径改为 `testing/{surface}/core.md`，增加旧结构检测 |

## Acceptance Criteria
- [ ] SKILL.md Convention 加载路径改为 `testing/{surface}/core.md` surface 目录遍历
- [ ] 检测到旧结构 convention 文件时输出迁移提示而非静默失败
- [ ] 生成的测试代码使用 per-surface build tag 命名（如 `cli_functional` 而非 `e2e`）

## Hard Rules
- 必须先加载 `docs/conventions/forge-distribution.md` 了解 plugin 分发模型
- `types/*.md` 为生成时的首要权威，core.md 断言偏好表源自 types/*.md

## Implementation Notes
- 策略信息权威关系：`types/*.md` 为 gen-test-scripts 生成时的首要权威；`core.md` 为 surface 策略权威。两者信息重叠部分以 `types/*.md` 为准
- 旧结构检测：检查 `docs/conventions/testing/` 下是否存在 `.md` 文件（非目录结构），若存在则提示运行 `/test-guide` 重新生成
