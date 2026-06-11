---
status: "completed"
started: "2026-05-26 22:40"
completed: "2026-05-26 22:43"
time_spent: "~3m"
---

# Task Record: 2 Clean documentation references to forge test commands

## Summary
Cleaned all documentation references to forge test commands (promote, run-journey, verify) across 12 files. Replaced CLI command references with /run-tests skill or tag-based promotion descriptions.

## Changes

### Files Created
无

### Files Modified
- README.md
- forge-cli/docs/OVERVIEW.md
- forge-cli/docs/OVERVIEW.zh.md
- docs/ARCHITECTURE.md
- docs/conventions/forge-cli-reference.md
- docs/conventions/forge-distribution.md
- docs/profile-authoring.md
- plugins/forge/skills/run-tests/SKILL.md
- plugins/forge/skills/consolidate-specs/SKILL.md
- plugins/forge/skills/gen-contracts/rules/journey-contract-model.md
- plugins/forge/skills/gen-journeys/rules/journey-contract-model.md
- plugins/forge/commands/run-tasks.md

### Key Decisions
无

## Document Metrics
12 files modified, 27 references removed, zero residue outside docs/features/

## Referenced Documents
- docs/proposals/remove-forge-test-command/proposal.md

## Review Status
completed

## Acceptance Criteria
- [x] Full-text search for forge test promote, forge test run-journey, forge test verify returns zero results (excluding docs/features/ history docs)
- [x] No documentation file instructs users or agents to run forge test subcommands

## Notes
All remaining matches are in docs/features/ (historical design records) and docs/proposals/ (historical proposals), both excluded per Hard Rules.
