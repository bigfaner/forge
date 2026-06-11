---
id: "3"
title: "Update skill files with fix-type derivation rule"
priority: "P1"
estimated_time: "2h"
dependencies: [1]
type: "doc"
mainSession: false
---

# 3: Update skill files with fix-type derivation rule

## Description

Replace all hardcoded `--type coding.fix` in error-handling instructions with a category-based derivation rule. When an agent creates a fix task, it must extract the source task's `TYPE` from claim output and derive the correct fix type: `doc.fix` for doc/eval categories, `coding.fix` for coding/test/validation/gate categories.

This is the core behavioral change that prevents doc task failures from spawning irrelevant code-level fix tasks.

## Reference Files

- `plugins/forge/agents/task-executor.md`: Line 50 — `forge task add --type coding.fix` in fix-task creation (source: proposal.md#In-Scope)
- `plugins/forge/commands/execute-task.md`: Lines 41, 74, 121, 123 — 4 hardcoded `coding.fix` occurrences in error handling (source: proposal.md#In-Scope)
- `plugins/forge/commands/run-tasks.md`: Lines 65, 69, 89, 111, 114 — 5 hardcoded `coding.fix` occurrences in dispatcher error handling (source: proposal.md#In-Scope)
- `plugins/forge/skills/submit-task/SKILL.md`: Line 111 — fix task creation on submit failure (source: proposal.md#In-Scope)
- `plugins/forge/skills/breakdown-tasks/SKILL.md`: Line 214 — Step 4 fix task template example (source: proposal.md#In-Scope)
- `plugins/forge/skills/quick-tasks/SKILL.md`: Line 228 — Step 4 fix task template example (source: proposal.md#In-Scope)
- `plugins/forge/skills/submit-task/data/record-format-coding.md`: Line 3 — Lists valid types, needs `doc.fix` added (source: proposal.md#In-Scope)

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/agents/task-executor.md` | Replace hardcoded `coding.fix` with derivation rule |
| `plugins/forge/commands/execute-task.md` | Replace 4 hardcoded `coding.fix` with derivation rule |
| `plugins/forge/commands/run-tasks.md` | Replace 5 hardcoded `coding.fix` with derivation rule |
| `plugins/forge/skills/submit-task/SKILL.md` | Replace hardcoded `coding.fix` in error handling; keep type reclassification mentions |
| `plugins/forge/skills/breakdown-tasks/SKILL.md` | Update Step 4 example to show derivation rule |
| `plugins/forge/skills/quick-tasks/SKILL.md` | Update Step 4 example to show derivation rule |
| `plugins/forge/skills/submit-task/data/record-format-coding.md` | Add `doc.fix` to valid types list |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria

- [ ] Error-handling instructions in task-executor.md, execute-task.md, run-tasks.md, submit-task/SKILL.md use derivation rule: extract `TASK_CATEGORY` from claim output, map doc/eval → `doc.fix`, coding/test/validation/gate → `coding.fix`
- [ ] Derivation rule table documented in at least one canonical location (run-tasks.md or execute-task.md) for agent reference
- [ ] `TYPE` and `TASK_CATEGORY` documented as extractable fields from `forge task claim` output in skill files
- [ ] `grep -rn "type coding\.fix" plugins/forge/ --include="*.md"` returns zero matches in error-handling contexts (informational mentions of `coding.fix` as a valid type or in reclassification examples are acceptable)

## Hard Rules

仅修改以下 7 个文件：plugins/forge/agents/task-executor.md, plugins/forge/commands/execute-task.md, plugins/forge/commands/run-tasks.md, plugins/forge/skills/submit-task/SKILL.md, plugins/forge/skills/breakdown-tasks/SKILL.md, plugins/forge/skills/quick-tasks/SKILL.md, plugins/forge/skills/submit-task/data/record-format-coding.md

## Implementation Notes

### Derivation rule

| Source Task Category | Fix Task Type | Rationale |
|----------------------|---------------|-----------|
| `doc`                | `doc.fix`     | 纯文档操作，修复也是改 .md 文件 |
| `eval`               | `doc.fix`     | 评估/修复的都是 .md spec 文件，无需改代码 |
| `coding`             | `coding.fix`  | 改代码 |
| `test`               | `coding.fix`  | 测试失败需改代码 |
| `validation`         | `coding.fix`  | 验证失败需改代码 |
| `gate`               | `coding.fix`  | 门禁检查含编译/单元测试，失败需改代码 |

### Distinction between error-handling and informational

Not all `coding.fix` mentions need replacement:
- **Error-handling** (MUST replace): `forge task add --type coding.fix` in catch-all error paths (timeout, blocked, missing instructions)
- **Informational** (keep as-is): type reclassification tables, valid type lists, "coding.fix" mentioned as a possible type value

### Agent behavior change

Agents currently have no mechanism to extract TYPE from claim output. The updated instructions must:
1. Tell agents to read `TASK_CATEGORY` from `forge task claim` output
2. Provide the derivation table as a lookup
3. Use the derived type in `forge task add --type <derived-type>`
