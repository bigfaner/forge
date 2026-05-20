---
id: "1"
title: "Add post-loop artifact commit step to run-tasks"
priority: "P2"
estimated_time: "30m"
dependencies: []
type: "doc"
mainSession: false
---

# 1: Add post-loop artifact commit step to run-tasks

## Description

After run-tasks completes all tasks and finishes knowledge extraction, uncommitted artifacts remain in the working tree. Add a "Commit Remaining Artifacts" section to the run-tasks Post-Completion flow that detects and commits these leftovers automatically.

## Reference Files
- `docs/proposals/post-loop-artifact-commit/proposal.md` — Source proposal

## Affected Files

### Create
| File | Description |
|------|-------------|
| _(none)_ | |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/commands/run-tasks.md` | Add "Commit Remaining Artifacts" section after Knowledge Review |

### Delete
| File | Reason |
|------|--------|
| _(none)_ | |

## Acceptance Criteria

- [ ] After Knowledge Review section (Step 6: Write confirmed knowledge), a new "### Commit Remaining Artifacts" section exists
- [ ] The section unconditionally runs `git status --porcelain` to detect uncommitted files
- [ ] Filters results to feature-scope paths only: `docs/features/<slug>/` and knowledge directories (`docs/decisions/`, `docs/lessons/`, `docs/conventions/`, `docs/business-rules/`)
- [ ] If filtered results are non-empty: `git add` the matched paths + `git commit` with message `chore(<slug>): commit post-loop artifacts`
- [ ] If filtered results are empty: skip silently, no output
- [ ] Does NOT run when loop ended due to 3 consecutive failures (matches Knowledge Review guard)

## Hard Rules

- Must follow Conventional Commits format for the generated commit message
- Must NOT use `git add -A` or `git add .` — only explicit paths from the filter

## Implementation Notes

- The section should be placed after the Knowledge Review subsection but before the end of the Post-Completion section
- The `<slug>` in paths and commit message refers to the active feature slug (already resolved in Step 0 of run-tasks)
- The knowledge directories filter ensures only project-level knowledge files are committed, not arbitrary docs changes
