---
feature: "worktree-start-idempotent"
created: "2026-06-09"
status: tasks
mode: quick
---

# Feature (Quick): worktree-start-idempotent

<!-- Status flow: tasks -> in-progress -> completed -->

## Documents

| Document | Path |
|----------|------|
| Proposal | ../../proposals/worktree-start-idempotent/proposal.md |

## Tasks

| ID | Title | Status | File |
|----|-------|--------|------|
| 1 | Rename `copy-files` → `includes` across codebase | pending | tasks/1-rename-copy-files-to-includes.md |
| 2 | Make `forge worktree start` idempotent — core behavior | pending | tasks/2-idempotent-start-core.md |
| 3 | Handle `--source-branch`, `--no-launch`, `--interactive` for existing worktrees | pending | tasks/3-idempotent-start-flags.md |
| T-test-1 | Test: Rename `copy-files` → `includes` across codebase | pending | tasks/T-test-1.md |
| T-test-2 | Test: Make `forge worktree start` idempotent — core behavior | pending | tasks/T-test-2.md |
| T-test-3 | Test: Handle `--source-branch`, `--no-launch`, `--interactive` for existing worktrees | pending | tasks/T-test-3.md |
| T-validate-all | Validate all tasks pass quality gate | pending | tasks/T-validate-all.md |
| T-clean-all | Cleanup completed tasks | pending | tasks/T-clean-all.md |
