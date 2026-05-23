---
status: "completed"
started: "2026-05-23 11:20"
completed: "2026-05-23 11:25"
time_spent: "~5m"
---

# Task Record: 3 prompt.go scope resolution 与 coverage 语言修复

## Summary
Fixed 3 issues in prompt.go: (1) resolveCoverage() now returns English text instead of Chinese, (2) coding.cleanup and coding.refactor types always use 'maintain' strategy regardless of frontmatter coverage field, preventing contradiction with 'no new tests' template directive, (3) added resolveScope() that checks project-type config to fall back scope when task scope mismatches project type (e.g. backend project + frontend scope -> empty scope). Updated existing tests to match new English output.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/prompt/prompt.go
- forge-cli/pkg/prompt/prompt_test.go

### Key Decisions
- cleanup/refactor types bypass frontmatter coverage override entirely - always use 'maintain' strategy with English text
- resolveScope reads project-type from .forge/config.yaml via forgeconfig.ReadConfig, single-scope projects (backend/frontend) fall back to empty scope on mismatch
- fullstack/mixed/library projects preserve scope as-is since they support multiple scopes

## Test Results
- **Tests Executed**: Yes
- **Passed**: 56
- **Failed**: 0
- **Coverage**: 90.0%

## Acceptance Criteria
- [x] scope resolution fallback: backend project + frontend scope -> default command (no scope suffix)
- [x] resolveCoverage() for coding.cleanup and coding.refactor does not inject percentage coverage directive
- [x] resolveCoverage() returns English text, no language mixing in templates
- [x] existing tests pass: prompt_test.go and scope_resolution_test.go
- [x] new tests cover scope fallback logic (backend project + frontend scope -> default command)

## Notes
无
