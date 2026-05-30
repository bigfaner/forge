---
status: "completed"
started: "2026-05-30 05:57"
completed: "2026-05-30 05:58"
time_spent: "~1m"
---

# Task Record: 7 Fix: init-justfile go.just test recipe uses Node.js command

## Summary
Replaced all Node.js/Playwright commands in go.just with Go ecosystem equivalents: test recipe now uses 'go test -v ./tests/...' (with '-run' for journey filtering), test-setup recipe now uses 'go mod download'

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/init-justfile/templates/go.just

### Key Decisions
无

## Document Metrics
2 recipes fixed (test, test-setup), 0 Node.js references remaining

## Referenced Documents
- docs/features/plugin-consistency-audit/reports/04-skills-batch-c.md
- docs/features/plugin-consistency-audit/reports/06-consolidated-report.md
- plugins/forge/skills/init-justfile/templates/node.just

## Review Status
final

## Acceptance Criteria
- [x] go.just test recipe no longer contains npx playwright test or any Node.js command
- [x] Replacement command uses Go ecosystem test command (go test ./...)
- [x] go.just has no remaining Node.js/Playwright references

## Notes
Server lifecycle and probe logic are language-agnostic and were preserved unchanged. test-setup simplified to 'go mod download' since Go does not need Playwright or npm.
