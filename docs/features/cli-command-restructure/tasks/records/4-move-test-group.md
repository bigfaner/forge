---
status: "completed"
started: "2026-05-23 01:57"
completed: "2026-05-23 02:15"
time_spent: "~18m"
---

# Task Record: 4 Move forge test group to test/ subdirectory

## Summary
Moved all forge test subcommand files from internal/cmd/ to internal/cmd/test/ subdirectory. Created Register() function for the test command group, mirroring the task sub-package pattern. Updated root.go to use testpkg.Cmd and testpkg.Register(). Tests split between test sub-package (unit tests for promote/verify/run-journey logic) and cmd package (integration tests using rootCmd). No circular dependencies: test sub-package imports only base/pkg packages, never internal/cmd.

## Changes

### Files Created
- forge-cli/internal/cmd/test/test.go
- forge-cli/internal/cmd/test/test_promote.go
- forge-cli/internal/cmd/test/test_verify.go
- forge-cli/internal/cmd/test/test_promote_test.go
- forge-cli/internal/cmd/test/test_verify_test.go
- forge-cli/internal/cmd/test/testmain_test.go

### Files Modified
- forge-cli/internal/cmd/root.go
- forge-cli/internal/cmd/root_test.go
- forge-cli/internal/cmd/test_test.go
- forge-cli/internal/cmd/test_verify_test.go

### Key Decisions
- Followed task sub-package pattern: exported Cmd var + Register() function
- Tests using rootCmd stayed in cmd package; tests using only Cmd moved to test sub-package
- journey_isolation_test.go merged into test/test_verify_test.go to consolidate command registration tests
- All error/output calls changed from cmd re-exports to direct base package imports

## Test Results
- **Tests Executed**: Yes
- **Passed**: 27
- **Failed**: 0
- **Coverage**: 52.8%

## Acceptance Criteria
- [x] All test subcommand files are in internal/cmd/test/
- [x] New package exports a Register() function
- [x] root.go updated to use the new package
- [x] go build ./... passes
- [x] go test ./... passes
- [x] forge test promote/run-journey/verify behavior unchanged

## Notes
Hard Rules verified: test sub-package does NOT import internal/cmd (no circular deps). All test-related unit tests moved to subdirectory. Integration tests using rootCmd remain in cmd package.
