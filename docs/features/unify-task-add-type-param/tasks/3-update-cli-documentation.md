---
id: "3"
title: "Update CLI documentation (WORKFLOW.md, OVERVIEW.md, README.md)"
priority: "P1"
estimated_time: "30min"
dependencies: ["1"]
type: "doc"
mainSession: false
---

# 3: Update CLI documentation (WORKFLOW.md, OVERVIEW.md, README.md)

## Description
Update all CLI documentation files that reference `--template` flag or template names. Replace with the new `--type` unified syntax. Ensure all examples, flag tables, and workflow descriptions reflect the removal of `--template`.

## Reference Files
- `proposal.md#Impact-Analysis` — lists affected documentation files with specific locations
- `proposal.md#Proposed-Solution` — defines the new syntax and behavior

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `forge-cli/docs/WORKFLOW.md` | Remove `--template` from flag table (line ~828); update all example commands using `--template fix-task`; update `addFixTask()` description; update IDPrefix example |
| `forge-cli/docs/OVERVIEW.md` | Line 312: update example command from `--template fix-task` to `--type coding.fix` |
| `README.md` | Line 154: replace `--template` row with `--type` in CLI parameter table |

## Acceptance Criteria
- [ ] No documentation file contains `--template` in the context of `forge task add`
- [ ] `WORKFLOW.md` flag table lists `--type` with description "Task type, auto-discovers matching template"
- [ ] All example commands use `--type coding.fix` instead of `--template fix-task`
- [ ] `README.md` CLI parameter table reflects new `--type` flag

## Hard Rules
- Only update `forge task add` related `--template` references — do NOT touch `task profile get --template` references (those are a different feature)
- `WORKFLOW.md` must sync with `add.go` per `forge-cli/CLAUDE.md` doc sync rules

## Implementation Notes
- `WORKFLOW.md` has the most changes: flag table, example commands, `addFixTask()` description, and IDPrefix mapping example
- Run `make check-docs` after changes to verify doc freshness
