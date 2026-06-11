---
status: "completed"
started: "2026-05-18 02:54"
completed: "2026-05-18 03:32"
time_spent: "~38m"
---

# Task Record: 8 CLI 重命名 + Tag-Based Promotion

## Summary
Renamed forge testing to forge test (breaking change). Removed graduate-tests skill entirely. Added forge test promote command with tag-based promotion (@feature -> @regression). Updated all skill files, documentation, and test files to reflect the rename.

## Changes

### Files Created
- forge-cli/internal/cmd/test.go
- forge-cli/internal/cmd/test_promote.go
- forge-cli/internal/cmd/test_promote_test.go
- forge-cli/internal/cmd/test_test.go

### Files Modified
- forge-cli/internal/cmd/root.go
- forge-cli/internal/cmd/root_test.go
- forge-cli/internal/cmd/test_verify.go
- forge-cli/internal/cmd/test_verify_test.go
- forge-cli/internal/cmd/journey_isolation_test.go
- forge-cli/internal/cmd/quality_gate.go
- forge-cli/scripts/version.txt
- forge-cli/docs/OVERVIEW.md
- forge-cli/docs/OVERVIEW.zh.md
- forge-cli/docs/WORKFLOW.md
- docs/conventions/forge-distribution.md
- plugins/forge/hooks/guide.md
- plugins/forge/references/shared/profile-detection.md
- plugins/forge/commands/run-tasks.md
- plugins/forge/skills/breakdown-tasks/SKILL.md
- plugins/forge/skills/breakdown-tasks/templates/graduate-tests.md
- plugins/forge/skills/breakdown-tasks/templates/validate-ux-task.md
- plugins/forge/skills/consolidate-specs/SKILL.md
- plugins/forge/skills/eval/SKILL.md
- plugins/forge/skills/eval/rubrics/test-cases.md
- plugins/forge/skills/gen-contracts/SKILL.md
- plugins/forge/skills/gen-test-cases/SKILL.md
- plugins/forge/skills/gen-test-scripts/SKILL.md
- plugins/forge/skills/init-justfile/SKILL.md
- plugins/forge/skills/quick-tasks/SKILL.md
- plugins/forge/skills/quick-tasks/templates/quick-graduate.md
- plugins/forge/skills/quick-tasks/templates/validate-ux-task.md
- plugins/forge/skills/run-e2e-tests/SKILL.md
- plugins/forge/skills/tech-design/SKILL.md

### Key Decisions
- Version bumped from 3.24.0 to 4.0.0 (major) since testing->test is a breaking CLI change
- graduate subcommand removed from test get -- replaced by promote at top level (forge test promote <journey>)
- promote command runs tests before promoting -- refuses promotion on failure (Hard Rule from AC)
- Tag replacement is language-aware: Go uses //go:build, Python uses @pytest.mark, Java uses @Tag(), Rust uses #[cfg]
- graduate-tests skill completely removed (Hard Rule: no compat shim)
- Contract files (_contracts/) are skipped during tag promotion to avoid modifying specs

## Test Results
- **Tests Executed**: Yes
- **Passed**: 612
- **Failed**: 0
- **Coverage**: 79.9%

## Acceptance Criteria
- [x] forge testing renamed to forge test, all subcommands work
- [x] forge test promote <journey> replaces @feature with @regression
- [x] promote confirms only tag changes via git diff
- [x] promote runs tests first, refuses on failure
- [x] CI tag filtering via --tags regression/feature
- [x] All existing tests pass, e2e-compile succeeds

## Notes
Breaking change: forge testing -> forge test. All 13 skill files updated. graduate-tests skill (3 template files) removed.
