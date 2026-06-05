---
status: "completed"
started: "2026-06-04 23:18"
completed: "2026-06-04 23:31"
time_spent: "~13m"
---

# Task Record: fix-2 Fix: Pre-existing test failures in qualitygate and scripts packages

## Summary
Fixed two pre-existing test failures: (1) TestCheckAllCompleted_NoProject failed because FindProjectRoot walks up from temp dir and finds package.json in user home -- replaced fragile subprocess isolation with direct function variable override; (2) install.sh and install.ps1 printed 'forge version' instead of 'forge --version' in verification instructions -- updated both scripts.

## Changes

### Files Created
无

### Files Modified
- forge-cli/internal/cmd/qualitygate/quality_gate.go
- forge-cli/internal/cmd/qualitygate/quality_gate_test.go
- forge-cli/scripts/install.sh
- forge-cli/scripts/install.ps1

### Key Decisions
- Introduced package-level findProjectRoot function variable (defaults to project.FindProjectRoot) instead of relying on subprocess env isolation for test reliability
- Fixed verification string in install scripts to match test expectations (forge --version)

## Test Results
- **Tests Executed**: Yes
- **Passed**: 78
- **Failed**: 0
- **Coverage**: 74.3%

## Acceptance Criteria
- [x] TestCheckAllCompleted_NoProject passes
- [x] TestInstallSh_PrintsVerificationInstructions passes
- [x] TestInstallPs1_PrintsVerificationInstructions passes

## Notes
Root cause 1: Windows dev env has package.json in user home, so FindProjectRoot walks up from temp dir and finds it. Root cause 2: install scripts used 'forge version' but tests expected 'forge --version'.
