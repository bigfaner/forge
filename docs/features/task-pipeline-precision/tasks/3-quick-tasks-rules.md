---
id: "3"
title: "Update quick-tasks split rules, complexity判定, Reference Files inline, remove task cap"
priority: "P0"
estimated_time: "1h"
dependencies: [1]
type: "doc"
mainSession: false
---

# 3: Update quick-tasks split rules, complexity判定, Reference Files inline, remove task cap

## Description

Update quick-tasks SKILL.md with all task generation improvements: split rules based on "independently verifiable" standard, complexity判定 with LLM override, Reference Files inline generation, and remove the 15 coding task cap.

## Reference Files
- `docs/proposals/task-pipeline-precision/proposal.md#Scope` — In Scope items: quick-tasks split rules, complexity判定, Reference Files inline, remove 15 cap
- `docs/proposals/task-pipeline-precision/proposal.md#Proposed-Solution` — defines the 3-layer precision control mechanism
- `docs/proposals/task-pipeline-precision/proposal.md#Constraints-&-Dependencies` — cleanTemplateOutput() conditional paragraph convention for task template compatibility

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/quick-tasks/SKILL.md` | Split rules, complexity判定, Reference Files generation, remove 15 cap |
| `plugins/forge/skills/quick-tasks/templates/task.md` | Add complexity field to frontmatter |
| `plugins/forge/commands/quick.md` | Remove 15 task cap from HARD-GATE |

## Acceptance Criteria

- [ ] Merge rule changed from "time estimation (<30min)" to "independently verifiable" standard
- [ ] AC max 6 rule added: "if a task has >6 AC, the scope is too large, split further"
- [ ] Multi-verb detection rule added: "task descriptions with connectors linking independent actions (rename + flatten + confirm) should be split by functional boundary"
- [ ] Complexity判定 logic: default heuristic (AC≤3 AND no Hard Rules AND Reference Files≤1 → low; AC>6 OR has Hard Rules → high; else → medium) + LLM judgment override guidance
- [ ] Reference Files generation changed from `proposal.md#Section-Title` pointers to inline precise info format (file path + specific change description)
- [ ] 15 coding task cap removed from SKILL.md HARD-GATE section
- [ ] `templates/task.md` frontmatter has `complexity: "{{COMPLEXITY}}"` field with default "medium"
- [ ] `quick.md` command's 15 task cap reference removed from HARD-GATE

## Hard Rules
{{HARD_RULES}}

## Implementation Notes

- The complexity判定 LLM override guidance should read: "如果静态指标与认知判断冲突（如 AC≤3 但涉及多文件架构变更），LLM 可根据认知判断覆盖默认 complexity 等级"
- The inline Reference Files format example: `- quality_gate.go: tests/e2e/results/raw-output.txt 路径需替换为 GetTestResultsDir()`
- When removing the 15 task cap from quick.md, also remove the ">15 coding tasks → STOP" logic from the Step 3→4 transition
