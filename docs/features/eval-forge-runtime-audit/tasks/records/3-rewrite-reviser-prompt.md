---
status: "completed"
started: "2026-05-15 21:51"
completed: "2026-05-15 21:53"
time_spent: "~2m"
---

# Task Record: 3 Rewrite reviser-prompt.md: two-layer fix strategy (safe-fix + guided-fix)

## Summary
Rewrote reviser-prompt.md with two-layer fix strategy: safe-fix (mechanical, no semantic change) for frontmatter, name mismatches, CLI flags, dead references, and status values; and guided-fix with 3 rules — Rule 1: instruction conflicts prefer guide.md as authority, Rule 2: content dedup with eval-loop-protocol extraction, Rule 3: bypass hardening with specific failure consequences (no empty prohibitions). Preserved HARD-RULE block and output format (FIXES APPLIED + FIXES SKIPPED).

## Changes

### Files Created
无

### Files Modified
- .claude/skills/eval-forge/templates/reviser-prompt.md

### Key Decisions
- Consolidated 12 fix categories from old reviser into 2 layers (5 safe-fix categories + 3 guided-fix rules)
- Removed command metadata fixes and plugin metadata fixes per task notes (D6 is only 50 pts, low priority)
- Removed safety marker fixes and model frontmatter fixes as they were not referenced in new rubric dimensions
- Added explicit forbidden patterns in Rule 3 to prevent empty prohibitions
- Added layer/rule traceability prefixes in output format for audit trail

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] Reviser prompt defines two fix layers: safe-fix and guided-fix
- [x] safe-fix covers: frontmatter missing fields, name-directory mismatch, CLI flag corrections, dead reference removal
- [x] guided-fix Rule 1: instruction conflicts, guide.md wins; SKILL.md changes to reference; if guide.md lacks concept, migrate most complete version
- [x] guided-fix Rule 2: content dedup, keep authoritative version; eval loop protocol extracts to references/shared/eval-loop-protocol.md
- [x] guided-fix Rule 3: bypass hardening, add minimal HARD-RULE with specific consequences, not empty prohibitions
- [x] HARD-RULE block preserved: only fix listed attack points, no refactoring, preserve existing intent
- [x] Output format: FIXES APPLIED + FIXES SKIPPED (requires human judgment)
- [x] Reviser reads audit report and rubric at specified paths

## Notes
Documentation-only task. Old reviser had 12 fix rule categories; new version consolidates into 2 layers aligned with the new 6-dimension rubric (D2/D3/D4).
