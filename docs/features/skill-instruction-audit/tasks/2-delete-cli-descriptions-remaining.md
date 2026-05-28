---
id: "2"
title: "Delete CLI behavior descriptions from remaining skills/commands"
priority: "P0"
estimated_time: "1.5h"
dependencies: [1]
type: "doc"
mainSession: false
---

# 2: Delete CLI behavior descriptions from remaining skills/commands

## Description

Remove CLI behavior descriptions from remaining skill/command files not covered in Task 1. Apply same three-category boundary rule.

## Reference Files
- `docs/proposals/skill-instruction-audit/proposal.md#CLI-描述删除边界规则`: Boundary rule table
- `plugins/forge/skills/gen-journeys/SKILL.md`: 4 forge surfaces example blocks (~30 lines)
- `plugins/forge/skills/run-tests/SKILL.md`: "segment prefix matching" implementation detail
- `plugins/forge/skills/eval/SKILL.md`: Repeated "Spawn as general-purpose agent" explanations
- `plugins/forge/skills/forensic/SKILL.md`: Outdated go build instructions
- `plugins/forge/skills/ui-design/SKILL.md`: Full bash script for config check
- `plugins/forge/commands/quick.md`: Behavioral descriptions of run-tasks and brainstorm

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/gen-journeys/SKILL.md` | Remove 4 forge surfaces example blocks; keep Exit Code contract table |
| `plugins/forge/skills/run-tests/SKILL.md` | Remove "segment prefix matching"; keep command |
| `plugins/forge/skills/eval/SKILL.md` | Remove repeated "Spawn as general-purpose agent via Agent tool" (2 occurrences) |
| `plugins/forge/skills/forensic/SKILL.md` | Replace outdated go build with "ensure forge CLI is installed" |
| `plugins/forge/skills/ui-design/SKILL.md` | Replace bash script with natural language instruction |
| `plugins/forge/commands/quick.md` | Remove "run-tasks reads index.json..." and "brainstorm runs its full interactive flow..." |

## Acceptance Criteria

- [ ] `gen-journeys/SKILL.md` has no example output blocks for forge surfaces; Exit Code table remains
- [ ] `run-tests/SKILL.md` has no "segment prefix matching"; command remains
- [ ] `eval/SKILL.md` has no repeated tool usage explanations
- [ ] `forensic/SKILL.md` has no go build or ~/.zcode-forge-cli references
- [ ] `ui-design/SKILL.md` config check is natural language, not bash script
- [ ] `quick.md` has no behavioral descriptions of run-tasks or brainstorm internals

## Hard Rules

- 仅修改上述 6 个文件

## Implementation Notes

- **gen-journeys**: 4 example blocks showing forge surfaces output are CLI docs, not agent instructions. Exit Code table provides decision logic.
- **eval**: Simplify "Spawn each composed prompt as a general-purpose agent via the Agent tool with model: 'sonnet'" to "Spawn scorer/reviser agent (model: sonnet)".
- **forensic**: ~/.zcode-forge-cli/task is outdated; current binary is `forge`. Source build assumes source access, contradicts distribution model.
- **ui-design**: Replace 12-line bash with: "Run `forge config get auto.eval.uiDesign`. true → auto-run; false → skip; unset → ask user."
- **quick.md**: Remove behavioral paragraphs at Step 1 (brainstorm internals) and Step 4 (run-tasks internals).
