---
created: 2026-05-16
author: "faner"
status: Approved
---

# Proposal: Proposal Status Lifecycle

## Problem

The `/quick` pipeline creates proposals with `status: Draft` but never updates this field after user confirmation or task completion, leaving 40+ proposals permanently stuck in Draft.

### Evidence

- 40+ existing proposals have `status: Draft` regardless of actual completion state
- A handful were manually changed to inconsistent values (`draft`, `proposed`, `approved`, `implemented`, `Superseded`, `proposal`)
- No automated status transition exists in any skill or CLI command
- `forge proposal list` and `forge feature status` display stale status values

### Urgency

As the number of proposals grows, the inability to distinguish active from completed proposals makes project tracking unreliable. This is a low-cost fix with high information-value return.

## Proposed Solution

Add automated proposal status transitions at two pipeline checkpoints in `/quick`:

1. **Draft → Approved**: when user confirms the proposal in Step 2
2. **Approved → Completed**: when all tasks finish in Step 4

Sync manifest status with proposal status, and update Go CLI display to recognize the new status values.

### Innovation Highlights

Straightforward lifecycle management — no novel approach needed. The key insight is embedding status updates at existing pipeline checkpoints rather than creating a separate status-management system.

## Requirements Analysis

### Key Scenarios

- User runs `/quick`, confirms proposal → status becomes Approved
- User runs `/quick`, aborts → status stays Draft
- `/run-tasks` completes all tasks → proposal status becomes Completed, manifest syncs
- `forge proposal list` shows Approved/Completed with appropriate labels

### Non-Functional Requirements

- Status updates must be atomic frontmatter edits (no content restructuring)
- Must not interfere with concurrent task execution

### Constraints & Dependencies

- Depends on existing `/quick` pipeline checkpoints (Step 2 confirmation, Step 4 run-tasks completion)
- Go CLI reads status for display only — no write-back from CLI

## Alternatives & Industry Benchmarking

### Industry Solutions

Standard issue/feature tracking systems (Jira, Linear, GitHub Issues) all have explicit status lifecycle with defined transitions.

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | No effort | Proposals remain stale, tracking useless | Rejected: defeats purpose of status field |
| Go CLI auto-infer only | — | No proposal.md changes | status field stays wrong, display-only fix | Rejected: addresses symptom not cause |
| **Markdown instructions + Go CLI** | Standard lifecycle | Proposals accurate end-to-end, CLI display matches | AI must reliably edit frontmatter | **Selected: end-to-end accuracy** |

## Feasibility Assessment

### Technical Feasibility

All changes are markdown instruction updates (3 skill files) and minor Go display logic. No architectural changes.

### Resource & Timeline

4-6 files to modify. Fits well within quick mode scope.

### Dependency Readiness

All dependent skills and CLI code exist and are stable.

## Scope

### In Scope

- `/quick` Step 2: add instruction to update proposal status to Approved after user confirms
- `/quick` Step 4: add instruction to update proposal status to Completed after all tasks finish
- Manifest status sync: when proposal → Completed, ensure manifest also reaches completed
- Go CLI: update `forge proposal list` and `forge feature status` to display Approved and Completed

### Out of Scope

- Full pipeline (`/write-prd` → `/tech-design`) status transitions
- Cleaning up 40+ existing proposals with inconsistent status values
- Go CLI status enum validation or enforcement
- New CLI commands for manual status management

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| AI fails to edit frontmatter correctly | L | M | Use explicit Edit instructions targeting the status line |
| Status update lost if run-tasks fails midway | M | L | Only update to Completed after all tasks confirmed done |
| Manifest and proposal status drift | M | M | Update both in same pipeline step |

## Success Criteria

- [ ] `/quick` user confirmation changes proposal.md status from Draft to Approved
- [ ] `/quick` all tasks completed changes proposal.md status from Approved to Completed
- [ ] manifest.md status syncs to completed when proposal reaches Completed
- [ ] `forge proposal list` displays Approved and Completed status values correctly
- [ ] `forge feature status` displays Approved and Completed status values correctly
- [ ] Abort at Step 2 leaves proposal status as Draft

## Next Steps

- Proceed to `/quick-tasks` to generate implementation tasks
