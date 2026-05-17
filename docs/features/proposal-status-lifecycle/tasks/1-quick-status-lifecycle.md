---
id: "1"
title: "Add proposal status lifecycle to /quick skill"
priority: "P1"
estimated_time: "30m"
dependencies: []
type: "documentation"
mainSession: false
---

# 1: Add proposal status lifecycle to /quick skill

## Description

The `/quick` pipeline creates proposals with `status: Draft` but never updates this field after user confirmation or task completion. Add automated status transitions at two existing pipeline checkpoints in `plugins/forge/commands/quick.md`:

1. **Step 2** (User Confirmation): After user selects "Yes, generate tasks", update `docs/proposals/<slug>/proposal.md` frontmatter `status` from `Draft` to `Approved`.
2. **Step 4** (Execute Tasks): After `/run-tasks` completes all tasks successfully, update both:
   - `docs/proposals/<slug>/proposal.md` frontmatter `status` from `Approved` to `Completed`
   - `docs/features/<slug>/manifest.md` frontmatter `status` from `tasks` to `completed`

## Reference Files
- `docs/proposals/proposal-status-lifecycle/proposal.md` — Source proposal

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/commands/quick.md` | Add status transition instructions at Step 2 and Step 4 |

## Acceptance Criteria
- [ ] `/quick` Step 2 instructions explicitly direct updating proposal.md frontmatter `status` from `Draft` to `Approved` when user confirms
- [ ] `/quick` Step 4 instructions explicitly direct updating proposal.md frontmatter `status` from `Approved` to `Completed` when all tasks finish
- [ ] `/quick` Step 4 instructions sync manifest.md frontmatter `status` to `completed` when proposal reaches Completed
- [ ] Abort at Step 2 leaves proposal status as Draft (no instruction to update on abort)

## Hard Rules
- Status updates MUST be atomic frontmatter edits only — use the Edit tool targeting only the `status:` line, not rewriting the entire file
- Both proposal.md and manifest.md status updates in Step 4 MUST happen together to prevent drift

## Implementation Notes
- The skill file is `plugins/forge/commands/quick.md`. Add instructions after the existing Step 2 confirmation block and after the Step 4 run-tasks invocation.
- Risk: AI fails to edit frontmatter correctly. Mitigation: use explicit Edit instructions targeting the `status:` line with enough surrounding context to be unique.
- Risk: Status update lost if run-tasks fails midway. Mitigation: only update to Completed after all tasks confirmed done.
