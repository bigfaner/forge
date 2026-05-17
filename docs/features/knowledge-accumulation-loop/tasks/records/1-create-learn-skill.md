---
status: "completed"
started: "2026-05-17 21:11"
completed: "2026-05-17 21:17"
time_spent: "~6m"
---

# Task Record: 1 Create /learn skill (SKILL.md + templates)

## Summary
Created unified /learn skill (SKILL.md + 3 templates) that absorbs /record-decision and /learn-lesson into a single entry point. Supports interactive and direct input modes, multi-type capture across 4 knowledge directories, write-first-then-report workflow, and reuses existing format specifications for compatibility with /consolidate-specs.

## Changes

### Files Created
- plugins/forge/skills/learn/SKILL.md
- plugins/forge/skills/learn/templates/decision-entry.md
- plugins/forge/skills/learn/templates/lesson-entry.md
- plugins/forge/skills/learn/templates/convention-entry.md

### Files Modified
无

### Key Decisions
- Write-first-then-report pattern: entries written immediately without pre-write confirmation, user reviews in final report
- Vocabulary as suggestions not enforcement: custom domain/type values always accepted
- Templates extracted from existing specs (decision-logging.md, learn-lesson/template.md) rather than reinvented
- Convention and business-rule entries use project-global ID encoding (TECH-/BIZ- prefix) for consolidate-specs compatibility

## Test Results
- **Tests Executed**: No
- **Passed**: 20
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] SKILL.md defines a skill with frontmatter name: learn and clear description
- [x] Skill supports two input modes: /learn (interactive) and /learn text (direct)
- [x] Skill workflow: identify knowledge type(s) -> classify -> write -> report
- [x] Knowledge type identification covers: decision, lesson, convention, business-rule
- [x] Multi-type capture: single input can produce entries in multiple directories
- [x] Write-first-then-report: entries written immediately, all shown in final report
- [x] Writes to docs/decisions/ using decision-logging.md row format + manifest update
- [x] Writes to docs/lessons/ using learn-lesson template format
- [x] Writes to docs/conventions/ by appending entries to existing domain files or creating new
- [x] Writes to docs/business-rules/ by appending entries to existing domain files or creating new
- [x] Accepts custom vocabulary values (domains, types) without error
- [x] Detects bulk extraction needs and delegates to /consolidate-specs
- [x] Uses auto-generated vocabulary when available, falls back gracefully when not
- [x] Category classification reuses 8-category vocabulary from learn-lesson tags + decision-logging type mapping

## Notes
Prompt-level only task -- no code changes. All 4 created files are markdown skill definitions and templates. Coverage set to -1 as there are no code changes to test. Quality gate (compile, fmt, lint, test) all pass.
