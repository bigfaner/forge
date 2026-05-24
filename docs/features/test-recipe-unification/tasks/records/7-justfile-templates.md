---
status: "completed"
started: "2026-05-24 21:19"
completed: "2026-05-24 21:31"
time_spent: "~12m"
---

# Task Record: 7 Overhaul justfile templates for two-layer test model

## Summary
Overhauled all 6 justfile templates (generic, go, node, python, rust, mixed) and root justfile for two-layer test model: renamed test to unit-test, removed e2e-test/e2e-setup/e2e-verify, added test (surface-level with optional journey param) and test-setup recipes with descriptive comments.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/init-justfile/templates/generic.just
- plugins/forge/skills/init-justfile/templates/go.just
- plugins/forge/skills/init-justfile/templates/node.just
- plugins/forge/skills/init-justfile/templates/python.just
- plugins/forge/skills/init-justfile/templates/rust.just
- plugins/forge/skills/init-justfile/templates/mixed.just
- justfile
- forge-cli/tests/justfile-integration/init_justfile_test.go
- forge-cli/tests/justfile-integration/forge_detection_test.go
- forge-cli/tests/justfile-integration/mixed_cli_test.go

### Key Decisions
- unit-test recipe uses language-appropriate commands: go test ./... (Go), npm test (Node), pytest (Python), cargo test (Rust)
- Go template unit-test omits -race flag per Implementation Notes (race detection controlled by CI)
- test recipe preserves server lifecycle management from old e2e-test but uses surface-agnostic wording in comments
- test-setup simplified from e2e-setup: removed force parameter, uses simpler idempotent check
- Root justfile test recipe dispatches to Go e2e tests (go test -tags=e2e) with journey filtering
- Removed test-e2e, e2e-setup, e2e-verify, e2e-compile, e2e-discover recipes from root justfile; replaced with test-setup and test-discover

## Test Results
- **Tests Executed**: Yes
- **Passed**: 30
- **Failed**: 0
- **Coverage**: 83.1%

## Acceptance Criteria
- [x] All 6 templates have unit-test recipe with language-appropriate command
- [x] All 6 templates have test recipe with optional journey parameter (test journey='')
- [x] All 6 templates have test-setup recipe
- [x] No e2e-test, e2e-setup, or e2e-verify recipes in any template
- [x] Comments: # unit-test: language-level unit tests and # test: surface-level advanced tests
- [x] Root justfile: test renamed to unit-test; new test recipe; ci updated
- [x] test recipe passes journey parameter to underlying test framework when non-empty

## Notes
Updated 3 test files (init_justfile_test.go, forge_detection_test.go, mixed_cli_test.go) to align assertions with new recipe names. E2e-tagged tests excluded from unit test run. Template tests reference e2e-test/e2e-verify in comments only (migration notes). Go source code (Tier 1), prompt templates (Tier 2), and docs (Tier 5) not touched — handled by other tasks.
