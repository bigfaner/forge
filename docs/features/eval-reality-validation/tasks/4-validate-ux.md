---
id: "4"
title: "validate-ux Implementation"
priority: "P1"
estimated_time: "4h"
dependencies: [1]
type: "documentation"
mainSession: false
---

# 4: validate-ux Implementation

## Description

Create the validate-ux eval type: new rubric (1000 pts, 10 dimensions), two-phase pre-processing logic in eval SKILL.md, CLI/Web/TUI project type detection via profile capabilities, and task templates for both breakdown-tasks and quick-tasks.

This is Batch 4 from the proposal — the highest-complexity batch.

## Reference Files
- `docs/proposals/eval-reality-validation/proposal.md` — Source proposal (contains full validate-ux design)
- `plugins/forge/skills/eval/SKILL.md` — Eval skill to extend
- `plugins/forge/skills/eval/rubrics/harness.md` — Reference for snapshot-based eval pattern

## Affected Files

### Create
| File | Description |
|------|-------------|
| `plugins/forge/skills/eval/rubrics/validate-ux.md` | Rubric with 10 dimensions, 1000 pts, 2板块 (UX Rules + PRD Flow) |
| `plugins/forge/skills/breakdown-tasks/templates/validate-ux-task.md` | Task template for full mode |
| `plugins/forge/skills/quick-tasks/templates/validate-ux-task.md` | Task template for quick mode |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/eval/SKILL.md` | Add validate-ux to all tables, add two-phase pre-processing logic, add ux-snapshot.md format spec |

### Delete
| File | Reason |
|------|--------|
| (none) | |

## Acceptance Criteria

- [ ] `validate-ux.md` rubric exists with `scale: 1000`, `target: 700`, `iterations: 1`, `type: validate-ux`, `context` frontmatter
- [ ] Rubric defines 10 dimensions: Error Actionability (120), Help Completeness (120), Output Clarity (90), Platform UX Rules (70), Flow Completeness (120), Output-Reality Consistency (120), Data & Side Effect (120), Idempotency & State Integrity (100), Cascade Effect (60), Friction Detection (80)
- [ ] Rubric `context` declares `conventions: [ux, cli, api]` and `business-rules: auto`
- [ ] Eval SKILL.md includes validate-ux pre-processing as two-phase: Phase 1 (main session: compile + run + collect ux-snapshot.md), Phase 2 (scorer: evaluate ux-snapshot.md)
- [ ] Project type detection uses `forge profile` capabilities: `cli` → CLI, `web-ui` → Web, `tui` → TUI
- [ ] Eval SKILL.md defines ux-snapshot.md format with Flow steps, Standalone Checks, and Effect Verification sections
- [ ] Eval SKILL.md defines PRD-to-operation translation strategies per project type (CLI/Web/TUI)
- [ ] Eval SKILL.md defines 7 operation impact verification types
- [ ] Task templates exist in both breakdown-tasks and quick-tasks positioned after T-test/T-quick steps
- [ ] validate-ux uses iterations=1, never triggers revise loop

## Hard Rules

- Do NOT modify doc-scorer.md or doc-reviser.md
- The rubric MUST total exactly 1000 points
- TUI pre-processing MUST be limited to non-interactive scenarios only (initial render, help, invalid input)
- validate-ux pre-processing MUST execute in a git worktree or temporary directory to avoid polluting project state
- `iterations: 1` is mandatory — no revise loop for validate-ux

## Implementation Notes

- The two-phase model follows the same pattern as `harness` type: pre-processing gathers state into a snapshot document, then scorer evaluates the snapshot.
- Phase 1 (main session) is the complex part:
  1. Read PRD → extract user flows
  2. Resolve project type via `forge profile` capabilities
  3. Compile and install the project binary
  4. For each flow: translate PRD actions to executable operations, run them, capture output
  5. Run standalone checks (--help, invalid command, --version)
  6. Execute effect verification (data diff, git diff, idempotency, state integrity)
  7. Write ux-snapshot.md
- Phase 2 is standard doc-scorer evaluation of ux-snapshot.md against the rubric.
- For CLI projects: operations are shell commands. For Web: use agent-browser with sitemap.json. For TUI: stdin pipe with non-interactive commands.
- Task templates should reference the PRD path, depend on all implementation + T-test tasks, and instruct the executor to run `forge eval --type validate-ux`.
