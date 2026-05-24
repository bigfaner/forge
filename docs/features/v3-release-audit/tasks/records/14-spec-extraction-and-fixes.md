---
status: "completed"
started: "2026-05-25 01:03"
completed: "2026-05-25 01:09"
time_spent: "~6m"
---

# Task Record: 14 Consolidate-specs extraction, UTF-8 handling, and CLI test path updates

## Summary
Extracted 78 lines from consolidate-specs SKILL.md to rules/drift-detection.md and rules/vocabulary-generation.md; assessed UTF-8 in guide.md (no changes needed); removed stale improve-harness references from CLI tests and skipped TC-003 referencing deleted templates/package.json.

## Changes

### Files Created
- plugins/forge/skills/consolidate-specs/rules/drift-detection.md
- plugins/forge/skills/consolidate-specs/rules/vocabulary-generation.md

### Files Modified
- plugins/forge/skills/consolidate-specs/SKILL.md
- forge-cli/tests/skill-ops/plugin_content_test.go
- forge-cli/tests/justfile-integration/mixed_cli_test.go
- forge-cli/tests/test-generation/gen_test_scripts_test.go

### Key Decisions
无

## Document Metrics
SKILL.md: 347->269 lines (-78); 2 new rules files; 3 test files fixed; guide.md UTF-8: no-op

## Referenced Documents
- plugins/forge/hooks/guide.md
- docs/conventions/forge-distribution.md

## Review Status
final

## Acceptance Criteria
- [x] consolidate-specs SKILL.md volume reduced >=50 lines
- [x] UTF-8 character handling assessed and impact recorded
- [x] CLI tests have no stale skill path references

## Notes
UTF-8 assessment was a no-op (guide.md non-ASCII chars are legitimate typographic symbols). improve-harness skill was replaced by run-tests; stale references removed. gen-test-scripts templates/ directory removed during profile v3 refactor; affected tests updated with t.Skip() or skipIfNoValidateScript().
