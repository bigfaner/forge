---
id: "13"
title: "Hook Unix parameter validation and misc file relocations"
priority: "P2"
estimated_time: "30m"
dependencies: ["6"]
type: "doc"
mainSession: false
---

# 13: Hook Unix parameter validation and misc file relocations

## Description
三个小修复：(1) Hook Unix 端参数校验检查；(2) validate-ux-pipeline.md 从 rubrics/ 移至 rules/；(3) 评估未暴露 skill 的 command 入口是否有必要保留。

## Reference Files
- `proposal.md#Scope` — P2.15: Hook Unix params; P2.16: validate-ux-pipeline relocation; P2.17: unexposed skill command eval

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/hooks/*.md` | Parameter validation review |

### Create
| File | Description |
|------|-------------|
| `plugins/forge/skills/*/rules/validate-ux-pipeline.md` | Relocated from rubrics/ |

### Delete
| File | Reason |
|------|--------|
| `plugins/forge/skills/*/rubrics/validate-ux-pipeline.md` | Relocated to rules/ |

## Acceptance Criteria
- [ ] Hook 参数格式已审查，问题已记录或修复
- [ ] validate-ux-pipeline.md 位于 rules/ 目录
- [ ] 未暴露 skill 的 command 入口已评估，决策已记录

## Hard Rules
- 移动文件后更新所有引用
- Hook 修改不影响运行时行为

## Implementation Notes
validate-ux-pipeline.md 的精确路径需确认。移动前 grep 确认引用位置。
