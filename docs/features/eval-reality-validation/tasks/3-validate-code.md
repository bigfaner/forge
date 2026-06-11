---
id: "3"
title: "validate-code Implementation"
priority: "P1"
estimated_time: "3h"
dependencies: [1]
type: "documentation"
mainSession: false
---

# 3: validate-code Implementation

## Description

Create the validate-code eval type: new rubric, pre-processing logic in eval SKILL.md, and task templates for both breakdown-tasks and quick-tasks. validate-code traces each PRD user scenario through git diff and implementation code to verify complete implementation paths.

This is Batch 3 from the proposal.

## Reference Files
- `docs/proposals/eval-reality-validation/proposal.md` — Source proposal
- `plugins/forge/skills/eval/rubrics/prd.md` — Reference rubric for format
- `plugins/forge/skills/eval/SKILL.md` — Eval skill to extend

## Affected Files

### Create
| File | Description |
|------|-------------|
| `plugins/forge/skills/eval/rubrics/validate-code.md` | Rubric for static code tracing evaluation |
| `plugins/forge/skills/breakdown-tasks/templates/validate-code-task.md` | Task template for full mode |
| `plugins/forge/skills/quick-tasks/templates/validate-code-task.md` | Task template for quick mode |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/eval/SKILL.md` | Add validate-code to Prerequisites table, Parameters table, Pre-Processing table, Report path, Rubric Reference table |

### Delete
| File | Reason |
|------|--------|
| (none) | |

## Acceptance Criteria

- [ ] `validate-code.md` rubric exists with `scale: 1000`, `target: 700`, `iterations: 1`, `type: validate-code`
- [ ] Rubric defines dimensions for: scenario traceability (PRD scenario → code path mapping), path completeness (full/partial/blocked per scenario), code-prd consistency
- [ ] Eval SKILL.md Prerequisites table includes `validate-code` entry requiring PRD + git diff
- [ ] Eval SKILL.md Pre-Processing table includes `validate-code` entry describing how to assemble PRD + git diff + code file list
- [ ] Eval SKILL.md Parameters table lists `validate-code` as a valid type
- [ ] Eval SKILL.md Rubric Reference table includes `validate-code` row with scale/target/iterations
- [ ] Task templates exist in both breakdown-tasks and quick-tasks with correct positioning (after implementation tasks, before T-test/T-quick steps)
- [ ] validate-code uses iterations=1, never triggers revise loop

## Hard Rules

- Do NOT modify doc-scorer.md or doc-reviser.md — prompt adaptation is done through eval SKILL.md prompt construction
- The rubric must use `iterations: 1` — validate-code produces a problem report, not a revised document
- Task templates must have `type: "gate"` and `breaking: false` — this is a verification step that produces a report, not a blocking gate

## Implementation Notes

- validate-code pre-processing: 1) Read PRD → extract user scenarios list. 2) Run `git diff <base-branch>...HEAD` to get changed files. 3) Compile changed file list. 4) Pass all to scorer as assembled input.
- The scorer prompt (constructed by eval SKILL.md) should instruct the scorer to: for each PRD scenario, trace through the git diff to find implementation evidence, then report "走通/走不通/部分实现" for each scenario.
- Report path: `docs/features/<slug>/eval/validate-code.md` (single report, no iterations).
- Task template content: the task instructs the executor to run `forge eval --type validate-code` on the current feature's PRD + implementation.
- Task template should reference the PRD path and depend on all implementation tasks.
