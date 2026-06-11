---
id: "6"
title: "Fix ambiguity and logic issues in pipeline commands/skills"
priority: "P1"
estimated_time: "1.5h"
dependencies: [5]
type: "doc"
mainSession: false
---

# 6: Fix ambiguity and logic issues in pipeline commands/skills

## Description

Fix ambiguous instructions, unclear preconditions, logic errors in pipeline-facing files. "歧义消除" + "逻辑修复" subcategories.

## Reference Files
- `plugins/forge/commands/quick.md`: Non-zero fallback = skip gate (should show gate)
- `plugins/forge/commands/execute-task.md`: MAIN_SESSION undefined, Step 1.5 indirect reference
- `plugins/forge/commands/run-tasks.md`: "successful cycle" boundary unclear, T-test-run undefined
- `plugins/forge/skills/gen-journeys/SKILL.md`: "Key Scenarios" matching rule unclear
- `plugins/forge/skills/submit-task/SKILL.md`: Type Reclassification trigger subjective
- `plugins/forge/commands/fix-bug.md`: $ARGUMENTS parsing unspecified, duplicate constraint

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/commands/quick.md` | Non-zero fallback: skip → show gate |
| `plugins/forge/commands/execute-task.md` | Define MAIN_SESSION upfront; make Step 1.5 self-contained |
| `plugins/forge/commands/run-tasks.md` | Define "successful" = STATUS completed; define T-test-run; add slug failure path |
| `plugins/forge/skills/gen-journeys/SKILL.md` | Add explicit section matching rule for quality degradation |
| `plugins/forge/skills/submit-task/SKILL.md` | Add objective criteria for type reclassification |
| `plugins/forge/commands/fix-bug.md` | Remove duplicate constraint; clarify $ARGUMENTS parsing |

## Acceptance Criteria

- [ ] `quick.md` Non-zero fallback = "show gate"
- [ ] `execute-task.md` defines MAIN_SESSION; Step 1.5 verify is self-contained
- [ ] `run-tasks.md` defines successful = STATUS completed; defines T-test-run; has slug failure path
- [ ] `gen-journeys` has case-insensitive "scenario" matching rule
- [ ] `submit-task` has objective type reclassification criteria
- [ ] `fix-bug` has no duplicate between E-I and HARD-GATE

## Hard Rules

- 仅修改上述 6 个文件
