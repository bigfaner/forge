---
id: "11"
title: "Fix: quick command Step 2 commit assumption"
priority: "P1"
estimated_time: "15min"
dependencies: []
type: "doc"
complexity: "low"
mainSession: false
---

# 11: Fix: quick command Step 2 commit assumption

## Description
`quick` command 的 Step 2 描述说 "The user already approved and committed the proposal in Step 1"，但 brainstorm skill 不自动 commit proposal。commit 需要在 brainstorm 完成后由 agent 执行。Step 2 的前提假设可能不成立。(Source: CMD-10, Report 05)

## Reference Files
- `docs/features/plugin-consistency-audit/reports/05-commands-agent-hooks.md#CMD-10`: P1 级 TIMING (source: Report 05)
- `plugins/forge/commands/quick.md`: 需修改的文件，Step 2 节 (source: audit finding)

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/commands/quick.md` | 修正 Step 2 的前提描述，删除 "and committed" 或在 Step 1-2 间增加显式 commit 步骤 |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] Step 2 的前提描述不再假设 proposal 已 committed
- [ ] 如果 brainstorm 已包含 commit 步骤（需验证），则仅删除多余的 "committed" 描述
- [ ] 如果 brainstorm 不含 commit 步骤，则在 Step 1-2 间增加 commit 指令

## Hard Rules
- 仅修改 `plugins/forge/commands/quick.md`
