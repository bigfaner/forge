---
id: "5"
title: "Slim generation domain (gen-sitemap + gen-journeys + gen-test-cases + gen-test-scripts)"
priority: "P1"
estimated_time: "1h"
dependencies: ["4"]
type: "doc"
mainSession: false
---

# 5: Slim generation domain (gen-sitemap + gen-journeys + gen-test-cases + gen-test-scripts)

## Description
对生成域的 4 个 skill 进行精简和消歧：gen-sitemap（229 行）、gen-journeys（211 行）、gen-test-cases（136 行）、gen-test-scripts（325 行）。

## Reference Files
- `docs/proposals/skill-slimming/proposal.md` — Splitting Heuristic
- `docs/conventions/skill-self-containment.md` — 自洽原则

## Affected Files

### Create
| File | Description |
|------|-------------|
| `plugins/forge/skills/gen-sitemap/rules/*.md` | gen-sitemap 规则细节 |
| `plugins/forge/skills/gen-journeys/rules/*.md` | gen-journeys 规则细节 |
| `plugins/forge/skills/gen-test-cases/rules/*.md` | gen-test-cases 规则细节 |
| `plugins/forge/skills/gen-test-scripts/rules/*.md` | gen-test-scripts 规则细节 |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/gen-sitemap/SKILL.md` | 保留流程骨架 |
| `plugins/forge/skills/gen-journeys/SKILL.md` | 保留流程骨架 |
| `plugins/forge/skills/gen-test-cases/SKILL.md` | 保留流程骨架 |
| `plugins/forge/skills/gen-test-scripts/SKILL.md` | 保留流程骨架 |

## Acceptance Criteria
- [ ] 每个 SKILL.md 行数 ≤ 350 行
- [ ] 所有步骤编号及描述保留
- [ ] 引用的辅助文件路径均存在可读
- [ ] 拆分风格与 Tier 1 保持一致

## Hard Rules
- 遵守 Splitting Heuristic
- 不改变 skill 的输入/输出契约

## Implementation Notes
- gen-sitemap 有 3 个辅助文件（182 行），gen-journeys 有 1 个（89 行）——评估现有文件是否已分担内容
- gen-test-cases 有 11 个辅助文件（843 行），gen-test-scripts 有 7 个（889 行）——已有丰富的 templates/ 和 types/ 目录，优先复用
- gen-test-cases 仅 136 行，可能只需要精简（消歧、清理冗余），不需要拆分
