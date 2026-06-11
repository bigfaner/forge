---
status: "completed"
started: "2026-06-09 22:37"
completed: "2026-06-09 22:42"
time_spent: "~5m"
---

# Task Record: 5 精简 SKILL.md 并删除废弃 rule 文件

## Summary
Simplified SKILL.md from 548 to 146 lines, deleted 6 obsolete rule files (server-lifecycle.md + 5 surface rules), updated agent workflow to use forge justfile scaffold CLI commands

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/init-justfile/SKILL.md

### Key Decisions
无

## Document Metrics
prompt layer: 180 lines (SKILL.md 146 + self-correction.md 34), down from 1645 (-89%)

## Referenced Documents
- docs/proposals/init-justfile-slim/proposal.md
- docs/conventions/forge-distribution.md

## Review Status
final

## Acceptance Criteria
- [x] SKILL.md + self-correction.md total lines <= 280
- [x] SKILL.md references forge justfile scaffold CLI command, workflow updated to Step 0-5
- [x] Delete 6 files: server-lifecycle.md, surfaces/{api,cli,mobile,tui,web}.md
- [x] rules/self-correction.md (34 lines) preserved unchanged
- [x] Convention Cold Start Fallback summarized in 5-10 lines

## Notes
Deleted files: rules/server-lifecycle.md, rules/surfaces/{api,cli,mobile,tui,web}.md, rules/surfaces/ directory. Removed: Step 1d Load Server Lifecycle Patterns, Step 3b Surface recipe generation details, Surface rule loading logic, Phase 1 Consistency Verification, Surface Gate Targets section, EXTREMELY-IMPORTANT duplicates, Notes duplicates.
