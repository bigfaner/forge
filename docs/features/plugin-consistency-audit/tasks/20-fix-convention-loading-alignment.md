---
id: "20"
title: "Fix: align Convention loading in 4 SKILL.md files"
priority: "P2"
estimated_time: "30min"
dependencies: []
type: "doc"
complexity: "medium"
mainSession: false
---

# 20: Fix: align Convention loading in 4 SKILL.md files

## Description
4 个 skill（gen-test-scripts, breakdown-tasks, tech-design, quick-tasks）的 SKILL.md Step 0 使用 `domains` frontmatter filtering 加载 Convention 文件，但 gen-test-scripts 的 `rules/convention-guide.md` 用 HARD-RULE 明确禁止此方式。v3.0.0 的 test profile 重构期间更新了 convention-guide.md，但 4 个 SKILL.md 未同步更新。

## Reference Files
- `docs/features/plugin-consistency-audit/reports/06-consolidated-report.md#Pattern-1`: Systemic pattern, Convention loading 不一致 (source: consolidated report)
- `plugins/forge/skills/gen-test-scripts/rules/convention-guide.md`: Convention 加载的权威规则 (source: audit finding)
- `plugins/forge/skills/{gen-test-scripts,breakdown-tasks,tech-design,quick-tasks}/SKILL.md`: 4 个需对齐的 SKILL.md Step 0

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/gen-test-scripts/SKILL.md` | 更新 Step 0 Convention 加载描述 |
| `plugins/forge/skills/breakdown-tasks/SKILL.md` | 更新 Step 0 Convention 加载描述 |
| `plugins/forge/skills/tech-design/SKILL.md` | 更新 Step 0 Convention 加载描述 |
| `plugins/forge/skills/quick-tasks/SKILL.md` | 更新 Step 0 Convention 加载描述 |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] 4 个 SKILL.md Step 0 的 Convention 加载方式与 `rules/convention-guide.md` 一致
- [ ] 如果 `domains` filtering 在某些 skill 中是合理的设计差异（如 breakdown-tasks 用于语言检测），则在 SKILL.md 中明确说明差异原因

## Hard Rules
- 仅修改上述 4 个 SKILL.md 文件的 Step 0 节

## Implementation Notes
- 先确认 convention-guide.md 的加载方式（index.md-based?），再决定统一方案
- 可能需要在 convention-guide.md 中增加说明：语言检测（breakdown-tasks/tech-design/quick-tasks）的 Convention 加载可以使用 domains，但测试生成（gen-test-scripts）不行
