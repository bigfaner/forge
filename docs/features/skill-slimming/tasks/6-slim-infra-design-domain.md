---
id: "6"
title: "Slim infra/design domain (init-justfile + ui-design + extract-design-md)"
priority: "P1"
estimated_time: "1h"
dependencies: ["5"]
type: "doc"
mainSession: false
---

# 6: Slim infra/design domain (init-justfile + ui-design + extract-design-md)

## Description
对基础设施/设计域的 3 个 skill 进行精简和拆分：init-justfile（387 行）、ui-design（314 行）、extract-design-md（242 行）。

## Reference Files
- `docs/proposals/skill-slimming/proposal.md` — Splitting Heuristic
- `docs/conventions/skill-self-containment.md` — 自洽原则

## Affected Files

### Create
| File | Description |
|------|-------------|
| `plugins/forge/skills/init-justfile/rules/*.md` | init-justfile 规则细节 |
| `plugins/forge/skills/ui-design/rules/*.md` | ui-design 规则细节 |
| `plugins/forge/skills/extract-design-md/rules/*.md` | extract-design-md 规则细节 |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/init-justfile/SKILL.md` | 保留流程骨架 |
| `plugins/forge/skills/ui-design/SKILL.md` | 保留流程骨架 |
| `plugins/forge/skills/extract-design-md/SKILL.md` | 保留流程骨架 |

## Acceptance Criteria
- [ ] 每个 SKILL.md 行数 ≤ 350 行
- [ ] 所有步骤编号及描述保留
- [ ] 引用的辅助文件路径均存在可读
- [ ] 拆分风格与 Tier 1 保持一致

## Hard Rules
- 遵守 Splitting Heuristic
- 不改变 skill 的输入/输出契约

## Implementation Notes
- init-justfile 有 6 个 templates/ 文件（1023 行）——大量 justfile 变体内联在 SKILL.md 或 templates/ 中，评估是否可进一步拆分
- ui-design 有 13 个辅助文件（1658 行）——已有 templates/platforms/ 和 templates/styles/，优先复用
- extract-design-md 有 3 个辅助文件（246 行），可能只需精简
