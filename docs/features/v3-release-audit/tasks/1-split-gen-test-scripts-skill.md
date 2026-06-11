---
id: "1"
title: "Split gen-test-scripts SKILL.md (527→≤350 lines)"
priority: "P0"
estimated_time: "1.5h"
dependencies: []
type: "doc"
mainSession: false
---

# 1: Split gen-test-scripts SKILL.md (527→≤350 lines)

## Description
gen-test-scripts SKILL.md 当前 527 行，超标 50%（上限 350 行）。提取 Step 0.5 和 Step 1 到独立 rules/ 文件，使 SKILL.md 回归约束。拆分后 SKILL.md 保留流程概述和引用加载，rules 文件承载详细逻辑。

## Reference Files
- `proposal.md#Proposed-Solution` — defines SKILL.md split as first P0 execution item
- `proposal.md#Scope` — P0.4 defines the split target (Step 0.5/1 extraction) and risk factors
- `proposal.md#Key-Risks` — SKILL.md split breaking reference chain, rollback via git stash

## Affected Files

### Create
| File | Description |
|------|-------------|
| `plugins/forge/skills/gen-test-scripts/rules/step-0.5-validation.md` | Step 0.5 validation logic extracted from SKILL.md |
| `plugins/forge/skills/gen-test-scripts/rules/step-1-contract-loading.md` | Step 1 contract loading logic extracted from SKILL.md |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/gen-test-scripts/SKILL.md` | Remove Step 0.5/1 inline content, add Load directives for new rules |

## Acceptance Criteria
- [ ] SKILL.md ≤ 350 行
- [ ] 新 rules/ 文件被 SKILL.md 通过 Load 引用（入度 ≥ 1）
- [ ] 拆分后 SKILL.md 流程完整，无断裂引用
- [ ] `wc -l plugins/forge/skills/gen-test-scripts/SKILL.md` ≤ 350

## Hard Rules
- 拆分前 `grep -r` 确认现有 rules 引用链不受影响
- 新 rules 文件需符合 skill-self-containment.md 规范
- 回滚方案：`git stash` 回归失败则降级为 P1

## Implementation Notes
风险：需显式 Load 引用，三路分支一致性，Mermaid 同步。gen-test-scripts 含 3 个现有 rules（convention-guide.md, quality-gates.md, run-to-learn.md），新增 rules 需保持命名一致。
