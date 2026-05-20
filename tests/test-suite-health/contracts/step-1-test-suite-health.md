---
step: 1
title: Test Suite Health Checks
journey: test-suite-health
---

# Step 1: Test Suite Health Checks

## Given
- Forge project source tree
- E2E test suite in tests/e2e/ and tests/*/ directories
- gen-journeys skill with SKILL.md and templates

## When
- Meta-tests check test suite quality invariants
- gen-journeys skill structure is verified

## Then
- Deleted test files and functions do not exist
- E2E test suite compiles with `go build -tags=e2e ./...`
- Zero unconditional t.Skip() calls
- Zero recursive `exec.Command("go", "test")` calls
- No static file text-grep tests remain
- No duplicate test files between root and features/ directories
- gen-journeys skill has valid structure, templates, and frontmatter
- gen-journeys SKILL.md references PRD, user stories, batch processing
- tui-ui-design directory is deleted
- TC-020 is removed from justfile-canonical-e2e

## Contract Dimensions
- **Actor**: CI pipeline or developer running quality checks on the test suite
- **Input**: Source tree file system, test files, skill directory
- **Output**: Pass/fail assertions on structural invariants
- **Error Cases**: any invariant violation -> test failure
- **Invariants**: tests must be self-contained, compilable, no unconditional skips
