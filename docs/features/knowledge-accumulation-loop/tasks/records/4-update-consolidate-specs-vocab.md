---
status: "completed"
started: "2026-05-17 21:30"
completed: "2026-05-17 21:45"
time_spent: "~15m"
---

# Task Record: 4 Update /consolidate-specs — add auto-vocabulary generation

## Summary
Added vocabulary generation step (Step 12) to /consolidate-specs SKILL.md. The new step scans all 4 knowledge directories (decisions, lessons, conventions, business-rules), aggregates types/domains/counts, includes the base 8-category vocabulary, and writes the index to docs/.vocabulary.md for /learn and auto-extract triggers. Existing steps 1-11 preserved unchanged; former Step 12 (Record Task) renumbered to Step 13. Added 10 e2e tests (TC-020 through TC-029) verifying step existence, scan coverage, base vocabulary, output structure, auto-generated marking, ordering, idempotency, and trigger references.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/consolidate-specs/SKILL.md
- forge-cli/tests/e2e/spec_drift_detection_cli_test.go

### Key Decisions
- Step 12 placed after Step 11 (commit spec changes) and before Step 13 (record task) so vocabulary reflects latest drift-fix state
- Vocabulary output written to docs/.vocabulary.md (dot-prefix indicates derived/auto-generated artifact)
- Vocabulary is suggestive not restrictive -- /learn and triggers accept values outside the vocabulary
- Base 8 categories always included regardless of knowledge directory contents

## Test Results
- **Tests Executed**: No
- **Passed**: 10
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] New step added after Step 11 drift fix, before Record Task
- [x] Vocabulary generation scans all 4 knowledge directories
- [x] Generates vocabulary index with types, domains, counts
- [x] Vocabulary output is prompt-readable format (markdown)
- [x] Step is idempotent -- regenerates on every run
- [x] /learn and auto-extract triggers can reference vocabulary
- [x] Works when knowledge directories are sparse or empty
- [x] Vocabulary includes base 8 categories even when no knowledge files exist
- [x] Generated vocabulary marked as auto-generated (not user-editable)
- [x] Existing consolidate-specs steps unchanged

## Notes
Prompt-level change only (SKILL.md). No Go source code modified. E2e tests verify SKILL.md structural correctness. Unit tests pass with >80% coverage across all packages.
