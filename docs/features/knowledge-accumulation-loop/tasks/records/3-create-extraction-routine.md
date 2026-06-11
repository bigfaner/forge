---
status: "completed"
started: "2026-05-17 21:18"
completed: "2026-05-17 21:23"
time_spent: "~5m"
---

# Task Record: 3 Create shared knowledge extraction routine

## Summary
Created shared knowledge extraction routine (knowledge-extraction.md) that defines the reusable prompt section for all 4 auto-extract trigger points. Includes parameterization by trigger context, artifact scanning scope per trigger type, 4 knowledge types with format references, extraction flow (scan-identify-extract-confirm-write), notable knowledge heuristics for all 4 types to achieve <30% false-positive rate, deduplication logic, and calling convention for trigger point inclusion.

## Changes

### Files Created
- plugins/forge/references/shared/knowledge-extraction.md

### Files Modified
无

### Key Decisions
- Structured heuristics as positive (NOTABLE) and negative (NOT notable) criteria pairs for each knowledge type to maximize conservative extraction
- Added deduplication step (Section 6) that checks existing knowledge dirs before presenting candidates
- Parameterized by trigger type and artifacts list so each trigger point specifies its own scanning scope
- Referenced existing format specs (decision-logging.md, learn-lesson template, consolidate-specs) rather than redefining formats

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] File defines a reusable prompt section that can be included by any trigger point
- [x] Defines the extraction flow: scan artifacts -> identify knowledge -> extract & summarize -> present for confirmation -> write on confirm
- [x] Knowledge identification covers 4 types: decisions, lessons, conventions, business rules
- [x] Defines heuristics for notable knowledge vs routine changes (to achieve <30% false-positive rate)
- [x] Silent when no notable knowledge detected — produces no output
- [x] Extracted knowledge presented for user confirmation before writing (AskUserQuestion)
- [x] Reuses same file formats as /learn skill (decision rows, lesson template, convention/business-rule entries)
- [x] Uses auto-generated vocabulary (from consolidate-specs) when available for classification suggestions
- [x] Parameterizable by trigger context: what artifacts to scan (varies by trigger point)
- [x] Defines the artifact scanning scope per trigger type (run-tasks, fix-bug, write-prd, tech-design)

## Notes
Reference/prompt file only — no executable code. Coverage auto-set to -1.0 since no testable code was produced. Quality gate (compile/fmt/lint/test) all passed.
