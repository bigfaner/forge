---
id: "1"
title: "Sync eval ecosystem data (H-1)"
priority: "P0"
estimated_time: "1h"
dependencies: []
type: "doc"
mainSession: false
---

# 1: Sync eval ecosystem data (H-1)

## Description

rubric-reference.md 中 journey 和 contract 的 scale/target 值与实际 rubric frontmatter 不一致；eval-journey/eval-contract 的 argument-hint 和 description 字段包含过时数据；eval/SKILL.md 的 scale 范围描述不准确。将 eval 生态的多个真相源同步为与 rubric frontmatter 一致。

## Reference Files
- `docs/proposals/forge-skill-audit/proposal.md` — H-1: rubric-reference.md 数据过时, Proposed Solution, Success Criteria, Regression Verification
- `plugins/forge/skills/eval/rules/rubric-reference.md`: Update journey/contract scale/target values (ref: H-1: rubric-reference.md 数据过时)
- `plugins/forge/commands/eval-journey.md`: Fix argument-hint and description (ref: H-1: rubric-reference.md 数据过时)
- `plugins/forge/commands/eval-contract.md`: Fix argument-hint and description (ref: H-1: rubric-reference.md 数据过时)
- `plugins/forge/skills/eval/SKILL.md`: Update scale range description (ref: H-1: rubric-reference.md 数据过时)

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/eval/rules/rubric-reference.md` | Update journey scale→1150/target→975; contract scale→1100/target→935; add maintenance comment at header |
| `plugins/forge/commands/eval-journey.md` | Update argument-hint --target to 975; description to 1150-point scale with all 7 dimensions including Workflow Coverage |
| `plugins/forge/commands/eval-contract.md` | Update argument-hint --target to 935; description to 1100-point scale with all 8 dimensions including Anchor Integrity and Fixture Specification |
| `plugins/forge/skills/eval/SKILL.md` | Remove "100-point" claim; update scale description to reflect actual scales |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] rubric-reference.md journey row: scale=1150, target=975 (matching rubric frontmatter)
- [ ] rubric-reference.md contract row: scale=1100, target=935 (matching rubric frontmatter)
- [ ] rubric-reference.md header has maintenance comment: this file is a secondary cache of rubric scale/target
- [ ] eval-journey.md argument-hint shows --target 975; description lists all 7 dimensions with 1150-point scale
- [ ] eval-contract.md argument-hint shows --target 935; description lists all 8 dimensions with 1100-point scale
- [ ] eval/SKILL.md description reflects actual supported scales (no "100-point" claim)

## Hard Rules
- Before modifying any plugin file, read `docs/conventions/forge-distribution.md`
- Only modify markdown files; no Go code changes
- Verify actual rubric frontmatter values before updating rubric-reference.md

## Implementation Notes
- eval 生态存在 5+ 处真相源，本任务仅同步前 4 处（rubric frontmatter、rubric-reference.md、命令 argument-hint、命令 description）；config 键默认值不在 scope 内
- 修复后运行回归验证：`grep "1000-point" plugins/forge/commands/eval-journey.md plugins/forge/commands/eval-contract.md` 确认无残留
