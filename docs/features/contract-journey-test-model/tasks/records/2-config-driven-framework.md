---
status: "completed"
started: "2026-05-18 00:13"
completed: "2026-05-18 00:47"
time_spent: "~34m"
---

# Task Record: 2 配置驱动框架系统

## Summary
Add config-driven test framework selection system to .forge/config.yaml. Projects can declare test-framework (e.g. mocha, pytest, go-testing) and test-command (e.g. go test ./...) to override language-based defaults. Built-in framework registry maps 6 frameworks to their code generation patterns. Zero-config falls back to language defaults. Fully backward compatible.

## Changes

### Files Created
- forge-cli/pkg/profile/framework.go
- forge-cli/pkg/profile/framework_test.go

### Files Modified
- forge-cli/pkg/profile/config.go
- forge-cli/pkg/profile/config_test.go
- forge-cli/internal/cmd/testing.go
- forge-cli/internal/cmd/testing_test.go
- forge-cli/internal/cmd/config_test.go
- forge-cli/internal/cmd/config_schema_test.go
- forge-cli/scripts/version.txt
- plugins/forge/references/shared/forge-config.schema.json
- plugins/forge/references/shared/forge-config.example.yaml

### Key Decisions
- Framework registry is a map[string]FrameworkInfo, not a language→framework hardcoded lookup -- language defaults exist but config override takes priority
- Custom framework names (not in builtins) return FrameworkInfo with only Name populated -- projects own the custom template resolution
- test-framework schema uses free-form string (not enum) to allow custom frameworks per Hard Rule 'do not hardcode language names into framework selection logic'
- ResolveTestFramework priority: config.TestFramework > DefaultFrameworkForLanguage(first language) > empty
- New forge testing framework subcommand exposes resolved framework info with PATTERN, FILES, SOURCE fields

## Test Results
- **Tests Executed**: Yes
- **Passed**: 17
- **Failed**: 0
- **Coverage**: 83.6%

## Acceptance Criteria
- [x] .forge/config.yaml supports test-framework and test-command fields
- [x] Declaring mocha generates describe/it pattern; pytest generates def test_*; go-testing generates func Test*
- [x] Zero-config uses built-in template defaults matching language defaults
- [x] Existing languages and interfaces fields continue working (backward compatible)

## Notes
17 new tests added: 11 in framework_test.go (GetBuiltinFramework, IsKnownFramework, KnownFrameworkNames, DefaultFrameworkForLanguage, ResolveTestFramework x7, RegisterCustomFramework x2), 4 in config_test.go (test-framework/test-command get value and absent), 2 in testing_test.go (framework auto-detect, config override, no-language). All 1254 existing tests pass with 0 failures.
