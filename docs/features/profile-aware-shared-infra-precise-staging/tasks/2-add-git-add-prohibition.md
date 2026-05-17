---
id: "2"
title: "Add git add -A prohibition to git-commit.md"
priority: "P0"
estimated_time: "15m"
dependencies: []
type: "documentation"
mainSession: false
---

# 2: Add git add -A prohibition to git-commit.md

## Description

git-commit.md instructs agents to "Stage appropriate files with `git add`" but never explicitly forbids broad staging commands. When agents can't determine their exact changed files (common during fix tasks), they fall back to `git add -A` which stages everything in the working directory — including unrelated untracked files and `.ts` residue.

## Reference Files
- `docs/proposals/profile-aware-shared-infra-precise-staging/proposal.md` — Source proposal
- `docs/lessons/gotcha-gen-test-scripts-ts-residue.md` — Incident showing 169-file commit from 2-file fix

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/commands/git-commit.md` | Add HARD-RULE prohibiting `git add -A`, `git add .`, `git add --all` |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria

- [ ] git-commit.md contains a new section (or expanded existing section) with a clear prohibition against `git add -A`, `git add .`, and `git add --all`
- [ ] The prohibition requires listing explicit file paths in the `git add` command
- [ ] The existing commit workflow (Steps 1-4, Task Completion Template) is otherwise unchanged

## Hard Rules

- Do NOT modify agent definitions (error-fixer.md, task-executor.md) — they delegate to git-commit skill
- The prohibition must be visually prominent (use `<HARD-RULE>` tag or equivalent emphasis)

## Implementation Notes

- Both error-fixer and task-executor delegate to git-commit via Skill tool, so this single fix covers all agent commits
- Place the prohibition before the Steps section or within the Steps section's staging step — wherever it's most visible to the agent
- Include a brief "Why" explaining the historical incident (169-file commit from 2-file fix) to give the agent context for compliance
