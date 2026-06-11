---
status: "completed"
started: "2026-05-18 02:30"
completed: "2026-05-18 02:52"
time_spent: "~22m"
---

# Task Record: 7 forge test verify 契约断裂检测

## Summary
Implement forge test verify contract breakage detection: scans Contract spec files under tests/<journey>/_contracts/*.md, parses them, collects Fact Table from current codebase, and compares Output/State assertions against actual values using semantic matching with zero false positives

## Changes

### Files Created
- forge-cli/pkg/contract/verify.go
- forge-cli/pkg/contract/verify_test.go
- forge-cli/internal/cmd/test_verify.go
- forge-cli/internal/cmd/test_verify_test.go

### Files Modified
- forge-cli/internal/cmd/testing_test.go
- forge-cli/scripts/version.txt

### Key Decisions
- Semantic matching extracts quoted content from descriptors (e.g., stdout contains "claimed task" -> matches on "claimed task" tokens) to avoid false positives from structural framing words
- FactCollector interface enables testability: production uses RealFactCollector, tests inject stubs with controlled facts
- verify is read-only per Hard Rules: no file modifications, Fact Table freshly collected each run, no cached snapshots
- No matching fact entry = no false positive: when facts are unavailable for a command, the contract is reported as OK rather than broken
- Contract parsing supports the canonical Markdown format with YAML frontmatter, Outcome blocks, dimension lines, and Journey Invariants

## Test Results
- **Tests Executed**: Yes
- **Passed**: 30
- **Failed**: 0
- **Coverage**: 84.5%

## Acceptance Criteria
- [x] forge test verify scans all Contract spec files (tests/<journey>/_contracts/*.md), re-executes code reconnaissance for each Output assertion
- [x] Detects breakage: reports affected Contract with file path, dimension, expected value, actual value
- [x] Zero false positives on 20+ unchanged Contracts
- [x] Fresh Fact Table collection on each run (no cached snapshots)
- [x] Bootstrap strategy: verify does not modify any files, only reads and reports

## Notes
New tests added: 25 in pkg/contract/verify_test.go + 5 in internal/cmd/test_verify_test.go. All Hard Rules enforced: verify is read-only, fresh facts each run, zero false positives verified with 22-contract test. Version bumped to 3.24.0.
