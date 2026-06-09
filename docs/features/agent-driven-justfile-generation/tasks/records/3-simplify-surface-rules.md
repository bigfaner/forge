---
status: "completed"
started: "2026-06-08 22:18"
completed: "2026-06-08 22:20"
time_spent: "~2m"
---

# Task Record: 3 Simplify 5 surface rule files to replace TODO stubs with Recipe Generation Requirements

## Summary
Simplified all 5 surface rule files (api/cli/tui/web/mobile) by replacing ## Recipe Template (Dual Platform) section with ## Recipe Generation Requirements section containing structural constraints for agent-driven recipe generation

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/init-justfile/rules/surfaces/api.md
- plugins/forge/skills/init-justfile/rules/surfaces/cli.md
- plugins/forge/skills/init-justfile/rules/surfaces/tui.md
- plugins/forge/skills/init-justfile/rules/surfaces/web.md
- plugins/forge/skills/init-justfile/rules/surfaces/mobile.md

### Key Decisions
无

## Document Metrics
5 files updated, ~130 lines of TODO stubs replaced with structured generation requirements per file

## Referenced Documents
- docs/proposals/agent-driven-justfile-generation/proposal.md

## Review Status
final

## Acceptance Criteria
- [x] All 5 surface rule files (api/cli/tui/web/mobile) updated
- [x] ## Recipe Template (Dual Platform) section replaced with ## Recipe Generation Requirements
- [x] Orchestration Sequence, Recipe Invocation Contract, Journey Filter Strategy sections preserved unchanged
- [x] Recipe Generation Requirements includes: naming rules, dual-platform attributes, user-customized marker, exit code semantics, test directory path rules
- [x] Dual-consumer consistency preserved: init-justfile and run-tests can both consume Recipe Invocation Contract

## Notes
Three surface categories identified with distinct generation requirements: (1) CLI/TUI - no server lifecycle, no aggregate recipe; (2) API/Web - full dev->probe->test->teardown with aggregate; (3) Mobile - extended sequence with test-setup step. All categories share common constraints: naming, dual-platform, exit codes, test paths, gate recipes.
