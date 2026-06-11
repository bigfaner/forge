---
id: "2"
title: "Create eval command wrappers for backward-compatible slash commands"
priority: "P0"
estimated_time: "1h"
dependencies: ["1"]
type: "documentation"
mainSession: false
---

# 2: Create eval command wrappers for backward-compatible slash commands

## Description
Create 7 thin command wrapper files in `plugins/forge/commands/` that delegate to the generic `eval` skill. This preserves all existing `/eval-*` slash commands so users and pipeline invocations see no change.

## Reference Files
- `docs/proposals/skill-rationalization/proposal.md` — Source proposal
- `plugins/forge/commands/quick.md` — Reference for command-to-skill delegation pattern

## Affected Files

### Create
| File | Description |
|------|-------------|
| `plugins/forge/commands/eval-proposal.md` | Wrapper → Skill("eval", "--type proposal") |
| `plugins/forge/commands/eval-prd.md` | Wrapper → Skill("eval", "--type prd") |
| `plugins/forge/commands/eval-design.md` | Wrapper → Skill("eval", "--type design") |
| `plugins/forge/commands/eval-ui.md` | Wrapper → Skill("eval", "--type ui") |
| `plugins/forge/commands/eval-test-cases.md` | Wrapper → Skill("eval", "--type test-cases") |
| `plugins/forge/commands/eval-consistency.md` | Wrapper → Skill("eval", "--type consistency") |
| `plugins/forge/commands/eval-harness.md` | Wrapper → Skill("eval", "--type harness") |

### Modify
| File | Changes |
|------|---------|
| (none) | |

### Delete
| File | Reason |
|------|--------|
| (none) | |

## Acceptance Criteria
- [ ] Each command wrapper has correct frontmatter (`name`, `description`)
- [ ] Each wrapper invokes `Skill("forge:eval", args="--type <type>")`
- [ ] `eval-harness` wrapper passes `--type harness` (not a separate skill invocation)
- [ ] Command descriptions match the original eval skill descriptions for system prompt consistency

## Hard Rules
- Command wrappers must be minimal — routing only, no eval logic
- Each wrapper file should be under 20 lines

## Implementation Notes
- Read `commands/quick.md` for the command-to-skill delegation pattern
- Copy descriptions from the original eval skill SKILL.md frontmatter
