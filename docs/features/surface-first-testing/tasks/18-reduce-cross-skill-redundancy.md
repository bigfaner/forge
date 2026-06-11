---
id: "18"
title: "Reduce cross-skill redundancy in surface detection + orchestration"
priority: "P2"
estimated_time: "1.5h"
dependencies: []
type: "doc"
mainSession: false
---

# 18: Reduce cross-skill redundancy in surface detection + orchestration

## Description
Surface Detection 流程在 gen-journeys 和 gen-contracts 中各内联约 30 行几乎完全相同的描述（CLI 命令、exit code 表、HARD-RULE）。编排序列在 init-justfile 和 run-tests 的 surface rule 文件中各维护一份。gen-contracts 内部 journey-contract-model.md 与 dimension-rules.md 的 Semantic Descriptors 规则几乎逐字复制。需要简化重复，改为引用 + 保留关键 HARD-RULE。

## Reference Files
- `docs/proposals/surface-first-testing/proposal.md` — Proposed Solution
- `plugins/forge/skills/gen-contracts/SKILL.md:56-80`: Surface Detection 内联描述
- `plugins/forge/skills/gen-journeys/SKILL.md:26-57`: Surface Detection 内联描述（几乎相同）
- `plugins/forge/skills/gen-contracts/rules/journey-contract-model.md:53-70`: Semantic Descriptors 重复
- `plugins/forge/skills/gen-contracts/rules/dimension-rules.md:24-45`: Semantic Descriptors 原始定义

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/gen-contracts/SKILL.md` | Surface Detection 小节简化为保留关键 HARD-RULE + "See gen-journeys SKILL.md Surface Detection for full flow" |
| `plugins/forge/skills/run-tests/SKILL.md` | Step 4 编排描述简化为引用 surface rule 文件，不重复内联失败处理 |
| `plugins/forge/skills/gen-contracts/rules/journey-contract-model.md` | Semantic Descriptors 规则简化为 "See rules/dimension-rules.md for full rules"，移除重复的正反例 |

## Acceptance Criteria
- [ ] gen-contracts SKILL.md Surface Detection 小节 ≤ 10 行（HARD-RULE + 引用）
- [ ] run-tests SKILL.md Step 4 编排描述不重复 surface rule 文件中已有的失败处理细节
- [ ] journey-contract-model.md Semantic Descriptors 段落简化为引用 dimension-rules.md
- [ ] 所有简化后的位置仍保留关键约束（HARD-RULE），不丢失信息

## Hard Rules
- 必须先加载 `docs/conventions/forge-distribution.md`
- 简化不能丢失任何 HARD-RULE 级别的约束
- 引用格式统一为 "See `<relative-path>` for details"

## Implementation Notes
- Surface Detection 的 HARD-RULE 是 "必须使用 forge surfaces CLI，禁止自行从目录结构推断"——这个约束必须在每个 skill 的 Surface Detection 处保留
- 编排序列的简化方向：SKILL.md 只描述高层流程（dev → probe → test → teardown），具体失败处理逻辑留给 surface rule 文件
