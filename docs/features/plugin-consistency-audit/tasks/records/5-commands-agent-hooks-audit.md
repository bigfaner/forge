---
status: "completed"
started: "2026-05-30 01:11"
completed: "2026-05-30 01:46"
time_spent: "~35m"
---

# Task Record: 5 Commands + Agent + Hooks deep audit

## Summary
Commands + Agent + Hooks deep audit: 18 commands, 1 agent (task-executor), and hooks/guide.md audited across Layer 1-3. Found 23 issues (4 P1, 15 P2, 4 P3). Key findings: fix-bug missing AskUserQuestion in allowed-tools and hardcoded Playwright; run-tasks MAIN_SESSION path missing submit-task/commit steps; quick assumes pre-committed proposal; guide.md omits mobile surface and many CLI commands.

## Changes

### Files Created
- docs/features/plugin-consistency-audit/reports/05-commands-agent-hooks.md

### Files Modified
无

### Key Decisions
无

## Document Metrics
23 findings across 3 layers: CONFLICT(4), INCOMPLETE(12), TIMING(2), REFERENCE(5). Severity: P0(0), P1(4), P2(15), P3(4). 11 of 18 commands confirmed clean.

## Referenced Documents
- docs/proposals/plugin-consistency-audit/proposal.md
- docs/features/plugin-consistency-audit/reports/01-inventory-structural.md
- plugins/forge/commands/clean-code.md
- plugins/forge/commands/eval-consistency.md
- plugins/forge/commands/eval-contract.md
- plugins/forge/commands/eval-design.md
- plugins/forge/commands/eval-journey.md
- plugins/forge/commands/eval-prd.md
- plugins/forge/commands/eval-proposal.md
- plugins/forge/commands/eval-ui.md
- plugins/forge/commands/execute-task.md
- plugins/forge/commands/extract-design-md.md
- plugins/forge/commands/fix-bug.md
- plugins/forge/commands/gen-sitemap.md
- plugins/forge/commands/git-checkout.md
- plugins/forge/commands/git-commit.md
- plugins/forge/commands/init-forge.md
- plugins/forge/commands/quick.md
- plugins/forge/commands/run-tasks.md
- plugins/forge/commands/simplify-skill.md
- plugins/forge/agents/task-executor.md
- plugins/forge/hooks/guide.md
- docs/reference/test-type-model.md

## Review Status
final

## Acceptance Criteria
- [x] All 18 command files read, internal flow step timing and reference consistency verified
- [x] task-executor agent file read, directive contradictions checked
- [x] hooks/guide.md read, script path existence verified (Layer 1 REFERENCE)
- [x] Script parameter description consistency verified (Layer 2 CONFLICT)
- [x] Each issue recorded with {component, file_path, layer, category, severity, description, fix_suggestion, confidence}

## Notes
guide.md does not reference hook scripts (confirmed by Task 1 report as intentional). guide.md is a standalone reference doc injected by session-start hook. 7 eval-* commands follow identical delegation pattern to forge:eval skill. Cross-component consistency verified for Fix-Type Derivation table across execute-task, run-tasks, and task-executor.
