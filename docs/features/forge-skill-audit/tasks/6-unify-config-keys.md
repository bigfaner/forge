---
id: "6"
title: "Unify auto.eval config key naming (M-1)"
priority: "P2"
estimated_time: "30m"
dependencies: [5]
type: "doc"
mainSession: false
---

# 6: Unify auto.eval config key naming (M-1)

## Description

auto.eval 配置键命名风格不一致：`auto.eval.proposal`/`auto.eval.prd` 使用全小写，`auto.eval.uiDesign`/`auto.eval.techDesign` 使用 camelCase。需统一为 kebab-case（`auto.eval.ui-design`, `auto.eval.tech-design`）。

## Reference Files
- `docs/proposals/forge-skill-audit/proposal.md` — M-1: auto.eval 配置键命名风格不一致, Proposed Solution, Risks
- `plugins/forge/skills/brainstorm/SKILL.md`: Rename auto.eval key (ref: M-1: auto.eval 配置键命名风格不一致)
- `plugins/forge/skills/write-prd/SKILL.md`: Rename auto.eval key (ref: M-1: auto.eval 配置键命名风格不一致)
- `plugins/forge/skills/ui-design/SKILL.md`: Rename auto.eval.uiDesign → auto.eval.ui-design (ref: M-1: auto.eval 配置键命名风格不一致)
- `plugins/forge/skills/tech-design/SKILL.md`: Rename auto.eval.techDesign → auto.eval.tech-design (ref: M-1: auto.eval 配置键命名风格不一致)

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/brainstorm/SKILL.md` | Verify auto.eval.proposal is already kebab-case |
| `plugins/forge/skills/write-prd/SKILL.md` | Verify auto.eval.prd is already kebab-case |
| `plugins/forge/skills/ui-design/SKILL.md` | Rename `auto.eval.uiDesign` → `auto.eval.ui-design` |
| `plugins/forge/skills/tech-design/SKILL.md` | Rename `auto.eval.techDesign` → `auto.eval.tech-design` |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] 所有 auto.eval 配置键在 skill markdown 中统一为 kebab-case（`auto.eval.ui-design`, `auto.eval.tech-design`）
- [ ] Implementation notes 明确标注 Go config reader alias 兼容需求为后续任务

## Hard Rules
- Before modifying any plugin file, read `docs/conventions/forge-distribution.md`
- Only modify markdown files; no Go code changes
- **先验证 Go config reader 是否支持 kebab-case 查询**：如果不支持，M-1 不可单独执行，仅在 markdown 侧标记 TODO 并创建跟踪 issue

## Implementation Notes
- Go config reader alias 兼容需要修改 Go 代码，超出本 proposal scope。如果 Go reader 不支持 kebab-case，用户按新 key 配置后 config 读取失败回退到默认值
- proposal Risks 表标记此为 "Likelihood: 中, Impact: 高"
