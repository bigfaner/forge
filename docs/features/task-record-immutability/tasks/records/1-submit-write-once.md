---
status: "completed"
started: "2026-05-17 01:17"
completed: "2026-05-17 01:32"
time_spent: "~15m"
---

# Task Record: 1 Submit write-once protection

## Summary
Add write-once protection to forge task submit: block record overwrite without --force, warn on stderr with --force, preserve normal path when record does not exist

## Changes

### Files Created
无

### Files Modified
- forge-cli/internal/cmd/submit.go
- forge-cli/internal/cmd/submit_test.go
- forge-cli/internal/cmd/integration_test.go
- forge-cli/scripts/version.txt

### Key Decisions
- Reused existing submitForce flag rather than adding a new flag
- Used os.Stat to detect existing record before os.WriteFile
- Used ErrValidation category via NewAIError per hard rules
- Warning message goes to stderr not stdout per hard rules

## Test Results
- **Tests Executed**: Yes
- **Passed**: 3
- **Failed**: 0
- **Coverage**: 80.7%

## Acceptance Criteria
- [x] forge task submit fails with descriptive error when record already exists (no --force)
- [x] forge task submit --force overwrites existing record with stderr warning
- [x] forge task submit succeeds normally when record does not exist
- [x] Existing submit tests continue to pass

## Notes
Also fixed TestForgeStateLifecycle which had Record field pointing to same path as task file (1.1.md), changed to records/1.1.md. Bumped version from 3.17.0 to 3.18.0 (minor: new feature).
