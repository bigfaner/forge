---
id: "4"
title: "Slim eval/quality domain (eval + gen-contracts + test-guide)"
priority: "P1"
estimated_time: "1h"
dependencies: ["3"]
type: "doc"
mainSession: false
---

# 4: Slim eval/quality domain (eval + gen-contracts + test-guide)

## Description
对评测/质量域的 3 个 skill 进行精简和拆分：eval（372 行）、gen-contracts（365 行）、test-guide（380 行）。按 Splitting Heuristic 规则处理每个文件。

## Reference Files
- `docs/proposals/skill-slimming/proposal.md` — Splitting Heuristic
- `docs/conventions/skill-self-containment.md` — 自洽原则

## Affected Files

### Create
| File | Description |
|------|-------------|
| `plugins/forge/skills/eval/rules/*.md` | eval 规则细节 |
| `plugins/forge/skills/gen-contracts/rules/*.md` | gen-contracts 规则细节 |
| `plugins/forge/skills/test-guide/rules/*.md` | test-guide 规则细节 |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/eval/SKILL.md` | 保留流程骨架 |
| `plugins/forge/skills/gen-contracts/SKILL.md` | 保留流程骨架 |
| `plugins/forge/skills/test-guide/SKILL.md` | 保留流程骨架 |

## Acceptance Criteria
- [ ] 每个 SKILL.md 行数 ≤ 350 行
- [ ] 所有步骤编号及描述保留
- [ ] 引用的辅助文件路径均存在可读
- [ ] 拆分风格与 Tier 1 保持一致

## Hard Rules
- 遵守 Splitting Heuristic
- 不改变 skill 的输入/输出契约
- eval 已有 28 个辅助文件（experts/、rubrics/）——不要破坏现有结构，只在需要时新增 rules/

## Implementation Notes
- eval 已有丰富的子目录结构，优先检查现有文件是否已承担 rules/ 的角色，避免重复拆分
- gen-contracts 有 2 个辅助文件（25 行），可能只需精简不需要拆分
- test-guide 有 1 个辅助文件（39 行），评估是否需要拆分
