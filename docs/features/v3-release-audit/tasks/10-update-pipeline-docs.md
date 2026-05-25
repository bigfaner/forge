---
id: "10"
title: "Update pipeline descriptions in guide.md and forge-distribution.md"
priority: "P1"
estimated_time: "30m"
dependencies: ["6"]
type: "doc"
mainSession: false
---

# 10: Update pipeline descriptions in guide.md and forge-distribution.md

## Description
guide.md 包含过时的 Pipeline 描述（v2 命名和流程），forge-distribution.md 的 Pipeline 描述与实际实现不一致。两者需与 v3.0.0 实际 Pipeline 流程对齐。

## Reference Files
- `proposal.md#Scope` — P1.10: guide.md updates; P1.11: forge-distribution.md alignment
- `proposal.md#Success-Criteria` — "guide.md Pipeline 描述与实现一致"; "forge-distribution.md Pipeline 描述与实现一致"

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `docs/conventions/guide.md` or `docs/reference/guide.md` | Update Pipeline descriptions to v3 naming and flow |
| `docs/conventions/forge-distribution.md` | Update Pipeline descriptions to match actual implementation |

## Acceptance Criteria
- [ ] guide.md Pipeline 描述与当前 skill 流程（/write-prd → /tech-design → /breakdown-tasks 等）一致
- [ ] forge-distribution.md Pipeline 描述与当前分发机制一致
- [ ] 无 v2 时代命令名残留

## Hard Rules
- 不重构文档结构，仅更新过时内容
- 先确认文件实际路径（可能在 docs/ 子目录）

## Implementation Notes
需先定位 guide.md 和 forge-distribution.md 的确切路径，然后与实际 skill/commands 目录结构比对。
