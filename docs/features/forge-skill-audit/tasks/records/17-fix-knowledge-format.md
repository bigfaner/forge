---
status: "completed"
started: "2026-06-10 21:18"
completed: "2026-06-10 21:19"
time_spent: "~1m"
---

# Task Record: 17 Fix knowledgeSave format description consistency (MINOR-H2)

## Summary
Unified auto.knowledgeSave output format description across 3 files to consistent wording: Parse the config output format quick:<val> full:<val> (e.g., quick:true full:false)

## Changes

### Files Created
无

### Files Modified
- plugins/forge/commands/fix-bug.md
- plugins/forge/skills/tech-design/rules/knowledge-extraction.md
- plugins/forge/skills/write-prd/rules/knowledge-extraction.md

### Key Decisions
无

## Document Metrics
3 files, 1 line each, format description unified

## Referenced Documents
- docs/conventions/forge-distribution.md

## Review Status
final

## Acceptance Criteria
- [x] 3 files use unified auto.knowledgeSave output format description wording

## Notes
fix-bug.md had 'plain text key:value pairs' while knowledge-extraction.md files used 'quick:<val> full:<val>'. Unified all to: Parse the config output format quick:<val> full:<val> (e.g., quick:true full:false)
