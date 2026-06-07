---
status: "completed"
started: "2026-06-07 23:54"
completed: "2026-06-07 23:55"
time_spent: "~1m"
---

# Task Record: 4 clean-code skill 增加 surface-aware gate recipe 指引

## Summary
Updated clean-code SKILL.md Step 3 Quality Gate with surface-aware prefixed recipe resolution logic: detect surface-key from task context, probe <key>-unit-test first, fallback to generic unit-test

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/clean-code/SKILL.md

### Key Decisions
无

## Document Metrics
Step 3 expanded from 27 lines to 43 lines, added Surface-Aware Recipe Resolution subsection with 3-step detection/probe/execute flow

## Referenced Documents
- docs/proposals/per-task-surface-scoped-gate/proposal.md
- docs/conventions/forge-distribution.md

## Review Status
final

## Acceptance Criteria
- [x] SKILL.md Step 3 指引 agent 检测当前任务是否有 surface-key，有则优先运行 just <key>-unit-test，不存在时回退到 just unit-test
- [x] 无 surface-key 时行为与改动前一致（运行 just unit-test，不存在则 skip）
- [x] 修改前必须先加载 docs/conventions/forge-distribution.md（CLAUDE.md MANDATORY 规则）

## Notes
Followed resolvePrefixedRecipe() fallback pattern from proposal.md. No surface-key → no probing, direct generic recipe execution. forge-distribution.md loaded before modification per CLAUDE.md MANDATORY rule.
