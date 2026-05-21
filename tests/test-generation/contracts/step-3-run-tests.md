---
step: 3
title: Run E2E Tests and Result Parsing
journey: test-generation
---

# Step 3: Run E2E Tests and Result Parsing

## Given
- A forge project with generated test scripts
- Convention files declaring Result Format (json-stream for Go)
- justfile with e2e-compile and e2e-test recipes

## When
- Tests are executed via the Convention-declared execution command
- `forge run-e2e-tests` is invoked

## Then
- Test results are parsed according to Convention-declared Result Format
- JSON-stream format: one JSON object per line with Action, Test, Package, Output fields
- TC IDs extracted from test function names via pattern matching
- Convention-based result format infrastructure is in place

## Contract Dimensions
- **Actor**: CLI user or skill executing test runs
- **Input**: Test files, Convention Result Format section, justfile recipes
- **Output**: Parsed test results with pass/fail/skip status, TC ID extraction
- **Side Effects**: Test result files written to results/ directory
- **Error Cases**: Compile failures, test failures, missing justfile recipes
- **Invariants**: Result format matches Convention declaration; TC ID extraction consistent with naming convention
