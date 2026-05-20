---
status: "completed"
started: "2026-05-20 14:17"
completed: "2026-05-20 14:27"
time_spent: "~10m"
---

# Task Record: 5 Slim generation domain (gen-sitemap + gen-journeys + gen-test-cases + gen-test-scripts)

## Summary
Slim generation domain: gen-sitemap (395->229 lines, extracted rules/schema.md, rules/page-exploration.md, rules/merge-validation.md), gen-test-scripts (350->325 lines, extracted rules/quality-gates.md). gen-journeys (211 lines) and gen-test-cases (136 lines) unchanged - already under 350 limit with no redundancy to trim.

## Changes

### Files Created
- plugins/forge/skills/gen-sitemap/rules/schema.md
- plugins/forge/skills/gen-sitemap/rules/page-exploration.md
- plugins/forge/skills/gen-sitemap/rules/merge-validation.md
- plugins/forge/skills/gen-test-scripts/rules/quality-gates.md

### Files Modified
- plugins/forge/skills/gen-sitemap/SKILL.md
- plugins/forge/skills/gen-test-scripts/SKILL.md

### Key Decisions
- gen-sitemap: extracted schema field table, element ID rules, dynamic route parameterization to rules/schema.md; element extraction rules and dynamic state exploration procedures to rules/page-exploration.md; merge/dedup/validation rules to rules/merge-validation.md
- gen-test-scripts: extracted antipattern guard table, duplicate name check, and error handling table to rules/quality-gates.md
- gen-journeys and gen-test-cases left unchanged: both well under 350-line limit, no redundant content to trim, no ambiguity items found

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] Each SKILL.md <= 350 lines
- [x] All step numbers and descriptions preserved
- [x] All referenced auxiliary file paths exist and are readable
- [x] Splitting style consistent with Tier 1

## Notes
无
