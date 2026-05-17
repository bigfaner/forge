---
id: "1"
title: "Cleanup stale record-task references"
priority: "P0"
estimated_time: "1h"
dependencies: []
scope: "all"
breaking: false
type: "cleanup"
mainSession: false
---

# 1: Cleanup stale record-task references

## Description

`record-task` was superseded by `submit-task` and deleted from the source tree, but 3 active files still reference it. Fix these stale references to prevent confusion and test failures.

## Reference Files
- `docs/proposals/skill-ecosystem-audit/proposal.md` — Source proposal, P2 finding #7

## Acceptance Criteria
- `grep -r 'record-task' plugins/forge/agents/ plugins/forge/skills/ plugins/forge/commands/` returns 0 hits (excluding eval/ reports and this proposal)
- `grep -r 'record-task' forge-cli/tests/e2e/` returns 0 hits or only references `submit-task`
- Go e2e test `justfile_mixed_cli_cli_test.go` passes after assertion update

## Hard Rules
- Do NOT modify historical docs (docs/lessons/, docs/decisions/) — those accurately reflect what happened at the time

## Implementation Notes
Three files to fix:
1. `plugins/forge/agents/task-executor.md:42` — Comment says "The submit-task skill internally calls record-task for metrics collection via `just test`." This is factually wrong. Fix to: "The submit-task skill collects metrics via `just test` before recording."
2. `plugins/forge/skills/breakdown-tasks/templates/consolidate-specs.md:88` — Says "Record task via `/record-task` skill." Change to `/submit-task`.
3. `forge-cli/tests/e2e/justfile_mixed_cli_cli_test.go:283-289` — Asserts "record-task" text exists in task-executor.md. Update assertion to check for "submit-task" instead.
