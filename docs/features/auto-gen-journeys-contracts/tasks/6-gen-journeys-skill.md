---
id: "6"
title: "gen-journeys SKILL.md 适配：proposal.md 输入 + AUTO_COMMIT 模式"
priority: "P0"
estimated_time: "1.5h"
dependencies: []
type: "doc"
mainSession: false
---

# 6: gen-journeys SKILL.md 适配：proposal.md 输入 + AUTO_COMMIT 模式

## Description

修改 gen-journeys SKILL.md 以支持：(1) proposal.md 作为替代输入源（Quick 模式没有 PRD user stories）；(2) AUTO_COMMIT=true 非交互模式（自动任务跳过用户审批）；(3) proposal.md 最低信息量检查。

## Reference Files
- `docs/proposals/auto-gen-journeys-contracts/proposal.md` — Source proposal
- `plugins/forge/skills/gen-journeys/SKILL.md` — 当前 SKILL.md

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/gen-journeys/SKILL.md` | 新增 proposal.md 输入路径、AUTO_COMMIT 条件行为、最低信息量检查 |

## Acceptance Criteria

- [ ] Prerequisites 部分明确区分两种输入模式：PRD 模式（prd-user-stories.md + prd-spec.md）和 Proposal 模式（proposal.md）
- [ ] Proposal 模式下，proposal.md 的 scope + success criteria 为硬性前提（缺失时任务 abort 并输出诊断信息）
- [ ] Proposal 模式下，key scenarios 缺失时降级为 smoke-level Journey（仅覆盖 happy path），并标注 quality=low
- [ ] Step 6 (Review & Commit) 新增 AUTO_COMMIT 条件路径：若任务上下文中包含 AUTO_COMMIT=true 指令，跳过用户审批，直接 git add + commit
- [ ] 手动调用 `/gen-journeys` 时行为不变（仍要求用户审批）
- [ ] Prerequisites 中 proposal.md 从 optional 升级为条件 required（当 PRD 文件不存在时为 required）

## Hard Rules

- 不修改核心 Journey 生成逻辑（Step 1-5 的识别工作流、风险分类、文件生成、验证流程）
- AUTO_COMMIT 路径必须同时执行 Step 5 的验证（Validate），仅跳过 Step 6 的人工审批

## Implementation Notes

- gen-journeys 当前已在 Prerequisites 中将 proposal.md 列为 optional。需要升级为条件 required：当 prd-user-stories.md 和 prd-spec.md 都不存在时，proposal.md 为 required
- proposal.md 的信息量检查：scope 和 success criteria 是 Journey 生成的最低要求，缺失会导致无法生成有意义的 Journey
