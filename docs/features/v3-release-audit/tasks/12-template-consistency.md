---
id: "12"
title: "Fix template variable naming and frontmatter consistency"
priority: "P2"
estimated_time: "30m"
dependencies: ["6"]
type: "doc"
mainSession: false
---

# 12: Fix template variable naming and frontmatter consistency

## Description
模板变量命名不统一（部分 `{{VAR}}` 部分 `${var}` 部分 `<var>`），模板 frontmatter 格式不一致。统一为规范格式。

## Reference Files
- `proposal.md#Scope` — P2.13: template variable naming unification; P2.14: frontmatter consistency

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/*/templates/*.md` | Standardize variable naming and frontmatter format |

## Acceptance Criteria
- [ ] 所有模板变量使用统一格式
- [ ] 模板 frontmatter 字段完整且格式一致

## Hard Rules
- 统一但不改变模板功能
- 不修改 skill 运行时行为

## Implementation Notes
需先扫描所有 templates/ 目录确定当前格式分布，再选择最常用格式作为标准。
