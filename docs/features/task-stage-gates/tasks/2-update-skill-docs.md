---
id: "2"
title: "Update skill docs to delegate gate/summary generation to forge task index"
priority: "P1"
estimated_time: "30min"
dependencies: ["1"]
scope: "all"
breaking: false
type: "implementation"
mainSession: false
---

# 2: Update skill docs to delegate gate/summary generation to forge task index

## Description
Update `/breakdown-tasks` and `/quick-tasks` skill documentation to reflect that `forge task index` now auto-generates `.summary` and `.gate` task files. Remove any manual gate/summary generation instructions from breakdown-tasks, and add stage-gate awareness to quick-tasks.

## Reference Files
- `docs/proposals/task-stage-gates/proposal.md` — Source proposal
- `plugins/forge/skills/breakdown-tasks/SKILL.md` — Currently generates gate/summary manually
- `plugins/forge/skills/quick-tasks/SKILL.md` — Currently has no gate/summary awareness
- `plugins/forge/skills/brainstorm/templates/proposal.md` — Proposal template (may need stage-gate mention)

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/breakdown-tasks/SKILL.md` | Remove manual gate/summary generation steps; add note that `forge task index` handles this automatically |
| `plugins/forge/skills/quick-tasks/SKILL.md` | Add note that `forge task index` auto-generates stage-gates for phases with >=2 business tasks |
| `forge-cli/scripts/version.txt` | Bump minor version (new feature) |

## Acceptance Criteria
- [ ] `/breakdown-tasks` SKILL.md no longer instructs the AI to manually write gate/summary task files
- [ ] `/breakdown-tasks` SKILL.md explicitly states that `forge task index` generates stage-gates
- [ ] `/quick-tasks` SKILL.md documents that `forge task index` generates stage-gates for qualifying phases
- [ ] No orphan references to manual gate/summary generation remain in either skill
- [ ] `forge-cli/scripts/version.txt` bumped (minor version for new feature)

## Implementation Notes
- The breakdown-tasks skill currently has a step that generates gate and summary task files from templates. This step should be replaced with a note that `forge task index` handles this.
- The quick-tasks skill just needs an informational note added — it already calls `forge task index` in Step 5.
