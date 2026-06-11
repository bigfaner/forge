---
status: "completed"
started: "2026-05-17 13:58"
completed: "2026-05-17 14:28"
time_spent: "~30m"
---

# Task Record: 2 Restructure profiles/ to languages/ with strategy files

## Summary
Restructure profiles/ to languages/ with language key directories. Renamed 6 profile dirs (go-test->go, web-playwright->javascript, rust-test->rust, pytest->python, java-junit->java, maestro->mobile). Removed all manifest.yaml files. Updated embed.go to use //go:embed all:languages. Replaced KnownProfiles with KnownLanguages. Updated languageCapabilities map to use new language keys matching D2 table. Updated detect.go to return language keys. Removed GetManifest, GetProfileInterfaces, manifest parsing code. Removed --manifest flag from CLI. Updated all callers and tests across profile, e2e, task, cmd packages.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/profile/embed.go
- forge-cli/pkg/profile/config.go
- forge-cli/pkg/profile/detect.go
- forge-cli/pkg/profile/embed_test.go
- forge-cli/pkg/profile/config_test.go
- forge-cli/pkg/profile/detect_test.go
- forge-cli/pkg/profile/autoconfig_test.go
- forge-cli/internal/cmd/profile.go
- forge-cli/internal/cmd/config.go
- forge-cli/internal/cmd/config_test.go
- forge-cli/internal/cmd/init.go
- forge-cli/internal/cmd/init_test.go
- forge-cli/pkg/e2e/e2e.go
- forge-cli/pkg/e2e/e2e_test.go
- forge-cli/pkg/e2e/actions_test.go
- forge-cli/pkg/task/types.go
- forge-cli/pkg/task/testgen.go
- forge-cli/pkg/task/testgen_test.go
- forge-cli/pkg/task/build_test.go
- forge-cli/pkg/task/autoconfig_test.go
- forge-cli/pkg/task/frontmatter_test.go
- forge-cli/pkg/just/just_test.go
- forge-cli/pkg/project/root_test.go
- forge-cli/tests/e2e/profile_cli_test.go
- forge-cli/tests/e2e/features/task-type-refinement/task_type_refinement_cli_test.go

### Key Decisions
- Kept DetectProfiles function name unchanged since it's called from many places; return values changed from profile names to language keys
- Removed manifest.yaml parsing entirely rather than migrating to a new format; capabilities are now purely in Go code
- Profile-to-language mapping per Hard Rules: go-test->go, web-playwright->javascript, rust-test->rust, pytest->python, java-junit->java, maestro->mobile
- go language capabilities changed from {api,cli,tui} to {api,cli} per proposal D2 table; web-playwright changed from {web-ui,api,cli} to {web-ui,api}
- Package stays at pkg/profile/ rather than moving to pkg/testing/ -- coordinate with Task 3/4 on final package path

## Test Results
- **Tests Executed**: Yes
- **Passed**: 20
- **Failed**: 0
- **Coverage**: 79.5%

## Acceptance Criteria
- [x] profiles/ directory no longer exists under forge-cli/pkg/profile/
- [x] languages/ directory contains 6 subdirectories: go, javascript, python, java, rust, mobile
- [x] Each language subdir has: generate.md, run.md, graduate.md, justfile-recipes, templates/
- [x] Zero manifest.yaml files remain in the codebase
- [x] embed.go uses //go:embed all:languages
- [x] languageCapabilities Go map defines capabilities for all 6 languages matching proposal D2 table
- [x] KnownProfiles constant/variable removed
- [x] go build ./... passes

## Notes
Strategy file content was not modified (Hard Rule). No hardcoded profiles/ references found in strategy files before rename. All 20 Go packages pass tests. Coverage at 79.5% for profile package.
