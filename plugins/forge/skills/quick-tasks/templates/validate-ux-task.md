---
id: "T-validate-ux"
title: "validate-ux: Dynamic UX Validation"
priority: "P1"
estimated_time: "2h"
dependencies: [{{ALL_QUICK_TASK_IDS}}]
breaking: false
type: "gate"
mainSession: true
---

# T-validate-ux: validate-ux: Dynamic UX Validation

## Description

Run dynamic UX validation against the PRD/proposal. Compile and install the project, execute user flows, capture outputs and effects into a ux-snapshot.md, then evaluate against the UX rubric. This produces a problem report, not a revised document.

Position: after T-quick-5 (last auto-generated test task), before all-completed hook.

## Instructions

### Step 1: Run validate-ux eval

```bash
forge eval --type validate-ux
```

This command:
1. Reads the PRD/proposal to extract user flows
2. Resolves project type via `forge test interfaces` (CLI/Web/TUI)
3. Compiles and installs the project binary
4. Translates user actions to executable operations per project type
5. Executes flows, captures outputs, runs effect verification
6. Writes ux-snapshot.md
7. Spawns doc-scorer to evaluate ux-snapshot.md against the rubric

### Step 2: Review the report

Read the generated report.

- If score >= 700 (target) -> proceed, note any friction points
- If score < 700 -> investigate failures and fix UX issues, then re-run

## Reference Files

- `docs/proposals/<slug>/proposal.md` -- Source proposal (primary input for quick mode)
- `docs/features/<slug>/prd/prd-spec.md` -- PRD with user flows (if exists)
- `plugins/forge/skills/eval/rubrics/validate-ux.md` -- Rubric definition (1000 pts, 10 dimensions)

## Acceptance Criteria

- [ ] `forge eval --type validate-ux` executed successfully
- [ ] ux-snapshot.md generated with Flow steps, Standalone Checks, and Effect Verification sections
- [ ] Every user flow appears in the snapshot with captured output
- [ ] Score report generated with per-dimension breakdown
- [ ] Any score below target is investigated and documented
- [ ] Record created via `/submit-task`

## Hard Rules

- MUST NOT modify implementation code -- this is verification only
- MUST execute in a git worktree or temporary directory to avoid polluting project state
- MUST NOT skip user flows -- every flow must appear in the snapshot
- For TUI projects: only non-interactive scenarios (initial render, help, invalid input) are in scope

## Implementation Notes

- The eval runs with `iterations: 1`, so there is no revise loop
- Pre-processing executes in a git worktree to isolate side effects
- Project type is auto-detected from `forge test interfaces`: `cli` -> CLI, `web-ui` -> Web, `tui` -> TUI
- CLI projects: operations are shell commands; Web: uses agent-browser with sitemap.json; TUI: stdin pipe with non-interactive commands
- In quick mode, the proposal's In Scope and Success Criteria sections serve as the flow source when no PRD exists
- If the feature is docs-only (no runtime code), this task can be skipped -- mark as completed with note "docs-only feature, no runtime to validate"
