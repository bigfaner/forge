---
id: "21"
title: "Fix: add intent-aware checks to write-prd + tech-design rule files"
priority: "P2"
estimated_time: "30min"
dependencies: []
type: "doc"
complexity: "medium"
mainSession: false
---

# 21: Fix: add intent-aware checks to write-prd + tech-design rule files

## Description
write-prd 和 tech-design 的 SKILL.md 正确实现了 `new-feature/refactor/cleanup` 三路 intent 分支，但其 rules 文件假设 `new-feature` intent，缺乏条件逻辑。具体：(1) write-prd `rules/self-check.md` 不检查 refactor/cleanup intent 下的必填字段；(2) tech-design `rules/design-quality-checks.md` 的 PRD coverage 和 DB schema 检查不区分 intent。

## Reference Files
- `docs/features/plugin-consistency-audit/reports/06-consolidated-report.md#Pattern-2`: Systemic pattern, Intent-aware 缺失 (source: consolidated report)
- `docs/features/plugin-consistency-audit/reports/02-skills-batch-a.md#WP-01`: write-prd self-check 缺少 intent 分支 (source: Report 02)
- `docs/features/plugin-consistency-audit/reports/02-skills-batch-a.md#TD-03`: tech-design quality checks 缺少 intent 分支 (source: Report 02)
- `docs/features/plugin-consistency-audit/reports/02-skills-batch-a.md#TD-04`: tech-design DB schema 检查不区分 intent (source: Report 02)
- `plugins/forge/skills/write-prd/rules/self-check.md`: 需添加 refactor/cleanup 条件检查
- `plugins/forge/skills/tech-design/rules/design-quality-checks.md`: 需添加 intent 条件

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/write-prd/rules/self-check.md` | 添加 refactor/cleanup intent 的检查条件 |
| `plugins/forge/skills/tech-design/rules/design-quality-checks.md` | 添加 intent-aware 条件（PRD coverage + DB schema） |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] `write-prd/rules/self-check.md` 包含 refactor/cleanup intent 的检查：验证 Change Scope、Constraints、Verification Criteria 三个必填字段；跳过 user stories、flow diagram、UI-related 检查
- [ ] `tech-design/rules/design-quality-checks.md` 的 PRD coverage 检查区分 intent：new-feature 从 user stories 取 AC，refactor/cleanup 从 Verification Criteria 取
- [ ] `tech-design/rules/design-quality-checks.md` 的 DB schema 检查在 refactor/cleanup intent 下跳过

## Hard Rules
- 仅修改上述 2 个 rules 文件
