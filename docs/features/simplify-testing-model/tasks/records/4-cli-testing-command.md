---
status: "completed"
started: "2026-05-17 14:40"
completed: "2026-05-17 15:00"
time_spent: "~20m"
---

# Task Record: 4 Implement forge testing CLI command group

## Summary
Replace forge profile command group with forge testing. New subcommands: detect (outputs detected languages), get generate/run/graduate/justfile/template (auto-detect language when --language flag not specified), interfaces (outputs interface types from config or detected language defaults). Removed profile.go entirely, created testing.go with cobra command hierarchy. Profile command is gone with no backward-compat aliases.

## Changes

### Files Created
- forge-cli/internal/cmd/testing.go
- forge-cli/internal/cmd/testing_test.go

### Files Modified
- forge-cli/internal/cmd/root.go
- forge-cli/internal/cmd/root_test.go

### Key Decisions
- Used cobra subcommand hierarchy: testing -> get -> generate/run/graduate/justfile/template instead of flags (--generate etc.) for cleaner CLI interface matching the proposal D4 spec
- Used PersistentFlags on testingGetCmd so --language flag is inherited by all get subcommands
- resolveLanguageFromFlags helper centralizes language resolution logic: --language flag > config override > auto-detect first language
- Output format uses same structured block pattern (--- delimiters, KEY: value) as existing commands for skill parsing stability

## Test Results
- **Tests Executed**: Yes
- **Passed**: 22
- **Failed**: 0
- **Coverage**: 79.9%

## Acceptance Criteria
- [x] forge testing detect outputs detected language(s)
- [x] forge testing get generate outputs generate.md content for auto-detected language
- [x] forge testing get generate --language go outputs Go strategy specifically
- [x] forge testing get run, get graduate, get justfile work for auto-detected language
- [x] forge testing get template <file> returns specified template file content
- [x] forge testing interfaces returns interface types (config.Interfaces > detected language defaults)
- [x] No language detected + no languages override -> exit code non-zero, stderr contains 'languages'
- [x] forge profile command removed entirely
- [x] go build ./... and go test ./... pass

## Notes
Multi-language project tests verify both detection of multiple languages and --language flag selection. Error case tested via resolveLanguageFromFlags helper unit test (os.Exit prevents direct integration test).
