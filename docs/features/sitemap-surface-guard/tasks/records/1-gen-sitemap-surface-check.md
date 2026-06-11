---
status: "completed"
started: "2026-06-01 21:21"
completed: "2026-06-01 21:23"
time_spent: "~2m"
---

# Task Record: 1 gen-web-sitemap: Add Step 0 surface type check

## Summary
Added Step 0 Surface Check to gen-web-sitemap SKILL.md Process Flow. Uses forge surfaces --json to detect web surface; stops with clear message if absent, passes through (including monorepo multi-surface) if present.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/gen-web-sitemap/SKILL.md

### Key Decisions
无

## Document Metrics
~25 lines added (Step 0 section + flow diagram update)

## Referenced Documents
- docs/proposals/sitemap-surface-guard/proposal.md
- plugins/forge/skills/gen-test-scripts/types/ui.md

## Review Status
final

## Acceptance Criteria
- [x] Step 0 executes forge surfaces --json and parses surface list
- [x] No web surface: STOP with clear message
- [x] Web surface present (incl. monorepo): PASS, proceed to Step 1
- [x] Empty output or command failure: treat as no web surface, abort

## Notes
Guard uses HARD-RULE tag for visibility. Wording follows ui.md Sitemap Resolution guard pattern.
