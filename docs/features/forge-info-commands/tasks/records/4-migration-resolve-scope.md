---
status: "completed"
started: "2026-05-14 16:41"
completed: "2026-05-14 16:50"
time_spent: "~9m"
---

# Task Record: 4 Migration: ResolveScope → config.yaml and doc updates

## Summary
Migrated ResolveScope() from subprocess-based `just project-type` to direct config.yaml read via profile.ReadConfig(). Updated all skill/hook docs referencing `just project-type` to use `forge config get project-type`. Removed `project-type` recipe from justfile and all 6 justfile templates (go, node, python, rust, generic, mixed).

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/just/just.go
- forge-cli/pkg/just/just_test.go
- plugins/forge/hooks/guide.md
- plugins/forge/skills/breakdown-tasks/SKILL.md
- plugins/forge/skills/init-justfile/SKILL.md
- plugins/forge/skills/init-justfile/templates/go.just
- plugins/forge/skills/init-justfile/templates/node.just
- plugins/forge/skills/init-justfile/templates/python.just
- plugins/forge/skills/init-justfile/templates/rust.just
- plugins/forge/skills/init-justfile/templates/generic.just
- plugins/forge/skills/init-justfile/templates/mixed.just
- justfile

### Key Decisions
- ResolveScope() imports profile.ReadConfig() directly instead of spawning subprocess - no subprocess call at all
- Empty/missing project-type returns empty string (backward-compatible skip scope behavior)
- Unknown project-type values still emit WARNING to stderr with the value for debugging

## Test Results
- **Tests Executed**: Yes
- **Passed**: 9
- **Failed**: 0
- **Coverage**: 98.1%

## Acceptance Criteria
- [x] ResolveScope() reads project-type from .forge/config.yaml via profile package, no subprocess call
- [x] All skill/hook docs use forge config get project-type instead of just project-type
- [x] justfile has no project-type: recipe
- [x] Skill justfile templates have no project-type: recipe
- [x] Existing ResolveScope() tests updated and passing
- [x] Scope resolution behavior is identical: mixed -> pass scope, frontend/backend -> skip scope, missing -> skip scope
- [x] Test coverage >= 80% for modified code

## Notes
Coverage at 98.1% for pkg/just. Hard rules satisfied: no subprocess spawn in ResolveScope(), backward-compatible behavior when config.yaml missing or project-type absent.
