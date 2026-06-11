---
feature: "proposal-status-lifecycle"
sources:
  - docs/proposals/proposal-status-lifecycle/proposal.md (quick mode — proposal serves as PRD)
generated: "2026-05-17"
---

# Test Cases: proposal-status-lifecycle

## Summary

| Type | Count |
|------|-------|
| UI   | 0    |
| TUI  | 0   |
| Mobile | 0  |
| API  | 0   |
| CLI  | 6   |
| **Total** | **6** |

---

## CLI Test Cases

## TC-001: Proposal status updates to Approved after quick confirmation
- **Source**: Proposal Success Criterion 1 — "`/quick` user confirmation changes proposal.md status from Draft to Approved"
- **Type**: CLI
- **Target**: cli/quick
- **Test ID**: cli/quick/proposal-status-updates-to-approved-after-quick-confirmation
- **Pre-conditions**: A proposal exists at `docs/proposals/<slug>/proposal.md` with `status: Draft` in frontmatter. The `/quick` pipeline has reached Step 2 (user confirmation).
- **Steps**:
  1. Run `/quick` and proceed to Step 2 confirmation
  2. Confirm the proposal (select "Yes, generate tasks")
  3. Read `docs/proposals/<slug>/proposal.md` frontmatter
- **Expected**: The `status` field in proposal.md frontmatter is `Approved` (was `Draft`)
- **Priority**: P0

## TC-002: Proposal status updates to Completed after all tasks finish
- **Source**: Proposal Success Criterion 2 — "`/quick` all tasks completed changes proposal.md status from Approved to Completed"
- **Type**: CLI
- **Target**: cli/quick
- **Test ID**: cli/quick/proposal-status-updates-to-completed-after-all-tasks-finish
- **Pre-conditions**: A proposal exists at `docs/proposals/<slug>/proposal.md` with `status: Approved` in frontmatter. All tasks for the feature have been generated and are ready to execute.
- **Steps**:
  1. Run `/run-tasks` for the feature until all tasks complete successfully
  2. Read `docs/proposals/<slug>/proposal.md` frontmatter
- **Expected**: The `status` field in proposal.md frontmatter is `Completed` (was `Approved`)
- **Priority**: P0

## TC-003: Manifest status syncs to completed when proposal reaches Completed
- **Source**: Proposal Success Criterion 3 — "manifest.md status syncs to completed when proposal reaches Completed"
- **Type**: CLI
- **Target**: cli/quick
- **Test ID**: cli/quick/manifest-status-syncs-to-completed-when-proposal-reaches-completed
- **Pre-conditions**: A feature manifest exists at `docs/features/<slug>/manifest.md` with `status: tasks`. A proposal at `docs/proposals/<slug>/proposal.md` has `status: Approved`. All tasks for the feature are ready to execute.
- **Steps**:
  1. Run `/run-tasks` for the feature until all tasks complete successfully
  2. Read `docs/features/<slug>/manifest.md` frontmatter
  3. Read `docs/proposals/<slug>/proposal.md` frontmatter
- **Expected**: manifest.md `status` field is `completed` and proposal.md `status` field is `Completed` — both updated atomically in the same pipeline step
- **Priority**: P0

## TC-004: Forge proposal list displays Approved status correctly
- **Source**: Proposal Success Criterion 4 — "`forge proposal list` displays Approved and Completed status values correctly"
- **Type**: CLI
- **Target**: cli/proposal-list
- **Test ID**: cli/proposal-list/displays-approved-status-correctly
- **Pre-conditions**: A proposal exists at `docs/proposals/<slug>/proposal.md` with `status: Approved` in frontmatter. The forge binary is built and available on PATH.
- **Steps**:
  1. Run `forge proposal list`
  2. Locate the row for the proposal with `status: Approved`
  3. Check the STATUS column value
- **Expected**: The STATUS column shows "Approved" for the proposal that has `status: Approved` in its frontmatter
- **Priority**: P0

## TC-005: Forge feature status displays Completed status correctly
- **Source**: Proposal Success Criterion 5 — "`forge feature status <slug>` correctly reflects when a feature's manifest status is completed"
- **Type**: CLI
- **Target**: cli/feature-status
- **Test ID**: cli/feature-status/displays-completed-status-correctly
- **Pre-conditions**: A feature exists at `docs/features/<slug>/manifest.md` with `status: completed` in frontmatter. The forge binary is built and available on PATH.
- **Steps**:
  1. Run `forge feature status <slug>`
  2. Check the output for status information
- **Expected**: The command output reflects "completed" status for the feature whose manifest has `status: completed`
- **Priority**: P1

## TC-006: Abort at Step 2 leaves proposal status as Draft
- **Source**: Proposal Success Criterion 6 — "Abort at Step 2 leaves proposal status as Draft"
- **Type**: CLI
- **Target**: cli/quick
- **Test ID**: cli/quick/abort-at-step-2-leaves-proposal-status-as-draft
- **Pre-conditions**: A proposal exists at `docs/proposals/<slug>/proposal.md` with `status: Draft` in frontmatter. The `/quick` pipeline has reached Step 2 (user confirmation).
- **Steps**:
  1. Run `/quick` and proceed to Step 2 confirmation
  2. Abort/reject the proposal (select "No" or abort)
  3. Read `docs/proposals/<slug>/proposal.md` frontmatter
- **Expected**: The `status` field in proposal.md frontmatter remains `Draft` (unchanged)
- **Priority**: P1

---

## Traceability

| TC ID | Source | Type | Target | Priority |
|-------|--------|------|--------|----------|
| TC-001 | Proposal Success Criterion 1 | CLI | cli/quick | P0 |
| TC-002 | Proposal Success Criterion 2 | CLI | cli/quick | P0 |
| TC-003 | Proposal Success Criterion 3 | CLI | cli/quick | P0 |
| TC-004 | Proposal Success Criterion 4 | CLI | cli/proposal-list | P0 |
| TC-005 | Proposal Success Criterion 5 | CLI | cli/feature-status | P1 |
| TC-006 | Proposal Success Criterion 6 | CLI | cli/quick | P1 |

---

## Route Validation

Discovered CLI commands (cobra registration):

| Command | Source |
|---------|--------|
| `forge proposal` | `forge-cli/internal/cmd/proposal.go:15` |
| `forge feature status` | `forge-cli/internal/cmd/feature.go:44` |

| TC ID | Target | Route/Command | Status |
|-------|--------|---------------|--------|
| TC-001 | cli/quick | `/quick` (skill pipeline) | Skill-invoked, not a CLI route |
| TC-002 | cli/quick | `/quick` (skill pipeline) | Skill-invoked, not a CLI route |
| TC-003 | cli/quick | `/quick` (skill pipeline) | Skill-invoked, not a CLI route |
| TC-004 | cli/proposal-list | `forge proposal list` | Matched: `forge proposal` at `proposal.go:15` (lists when no slug arg provided) |
| TC-005 | cli/feature-status | `forge feature status <slug>` | Matched: `forge feature status` at `feature.go:44` |
| TC-006 | cli/quick | `/quick` (skill pipeline) | Skill-invoked, not a CLI route |
