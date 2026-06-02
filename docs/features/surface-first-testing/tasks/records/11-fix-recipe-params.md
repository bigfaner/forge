---
status: "completed"
started: "2026-06-02 22:51"
completed: "2026-06-02 22:55"
time_spent: "~4m"
---

# Task Record: 11 Fix recipe parameter mismatch + run-tests surface-key usage

## Summary
Fixed recipe parameter mismatch between init-justfile and run-tests: (1) init-justfile SKILL.md and 5 surface rule files now declare <key>-test recipe with optional [journey] parameter; (2) run-tests SKILL.md Step 1 adds recipe-prefix logic (surface-key for multi-surface, surface-type for single-surface); (3) run-tests 5 surface rule files updated to use <recipe-prefix>-test <journey> in orchestration tables and per-journey pseudocode

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/init-justfile/SKILL.md
- plugins/forge/skills/init-justfile/rules/surfaces/cli.md
- plugins/forge/skills/init-justfile/rules/surfaces/api.md
- plugins/forge/skills/init-justfile/rules/surfaces/web.md
- plugins/forge/skills/init-justfile/rules/surfaces/tui.md
- plugins/forge/skills/init-justfile/rules/surfaces/mobile.md
- plugins/forge/skills/run-tests/SKILL.md
- plugins/forge/skills/run-tests/rules/surfaces/cli.md
- plugins/forge/skills/run-tests/rules/surfaces/api.md
- plugins/forge/skills/run-tests/rules/surfaces/web.md
- plugins/forge/skills/run-tests/rules/surfaces/tui.md
- plugins/forge/skills/run-tests/rules/surfaces/mobile.md

### Key Decisions
无

## Document Metrics
12 files modified, 5 AC items verified PASS

## Referenced Documents
- docs/proposals/surface-first-testing/proposal.md
- docs/conventions/forge-distribution.md

## Review Status
final

## Acceptance Criteria
- [x] init-justfile SKILL.md declares recipe signature as just <key>-test [journey]
- [x] init-justfile 5 surface rule files accept optional journey parameter
- [x] run-tests SKILL.md uses surface-key from forge surfaces --json for multi-surface projects
- [x] run-tests 5 surface rule files per-journey pseudocode aligned with init-justfile recipe signature
- [x] Single-surface projects use surface-type as recipe prefix (backward compatible)

## Notes
Recipe template syntax uses just's named optional parameter: `journey=''`. Run-tests introduces `recipe-prefix` concept that resolves to surface-key (multi-surface) or surface-type (single-surface), aligning with init-justfile's naming convention.
