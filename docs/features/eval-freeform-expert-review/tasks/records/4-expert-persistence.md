---
status: "completed"
started: "2026-05-23 20:07"
completed: "2026-05-23 20:09"
time_spent: "~2m"
---

# Task Record: 4 Expert Persistence, Reuse & Deprecation

## Summary
Created freeform-expert-persistence.md defining expert persistence directory structure, YAML front matter schema, reuse matching logic (Jaccard similarity + weighted scoring), quality tracking mechanism, auto-deprecation (3 consecutive non-substantive reviews), and manual deprecation rules.

## Changes

### Files Created
- plugins/forge/skills/eval/rules/freeform-expert-persistence.md

### Files Modified
无

### Key Decisions
无

## Document Metrics
1 file created, ~260 lines, 6 acceptance criteria met

## Referenced Documents
- docs/proposals/eval-freeform-expert-review/proposal.md
- plugins/forge/skills/eval/experts/freeform/expert-template.md
- plugins/forge/skills/eval/experts/freeform/expert-inference.md
- plugins/forge/skills/eval/rules/freeform-injection.md
- plugins/forge/skills/eval/rules/scorer-composition.md

## Review Status
final

## Acceptance Criteria
- [x] Define docs/experts/ directory structure with per-expert files named by domain slug
- [x] Define YAML front matter required fields: domain, background, review_style, generated_for, created_at, review_history, deprecated
- [x] Define reuse matching logic: load non-deprecated experts, extract keywords, compute Jaccard overlap, present best match via AskUserQuestion
- [x] Define quality tracking: record substantive_change in review_history (rubric_delta >= 15 or attack_points_changed)
- [x] Define auto-deprecation: 3 consecutive non-substantive reviews trigger deprecated: true
- [x] Define manual deprecation: user edits deprecated field in YAML front matter

## Notes
Reuse matching uses Jaccard similarity as primary metric with extended weighted scoring aligned to expert-inference.md. Persistence flow summary diagram included at end of document.
