---
id: "3"
title: "Fix tests/e2e/ feature tests to use TestMain-built binary"
priority: "P1"
estimated_time: "1-2h"
dependencies: []
scope: "backend"
breaking: false
type: "coding.enhancement"
mainSession: false
---

# 3: Fix tests/e2e/ feature tests to use TestMain-built binary

## Description

`tests/e2e/main_test.go` correctly builds an isolated binary into a temp dir and exposes it via the `forgeBinary` variable. However, several feature test files under `tests/e2e/features/` bypass this by calling `exec.Command("forge", ...)` directly, relying on PATH. Additionally, `tests/e2e/features/test-knowledge-convention-driven/test_helpers_test.go` has its own duplicate TestMain build pattern. Unify all tests to use the top-level `forgeBinary` variable.

## Reference Files
- `docs/proposals/e2e-isolated-build-artifact/proposal.md` — Source proposal
- `tests/e2e/main_test.go` — Existing TestMain build pattern

## Acceptance Criteria
- All test files under `tests/e2e/` use the `forgeBinary` variable from `main_test.go` instead of bare `"forge"`
- `grep -r 'exec.Command("forge"' tests/e2e/ --include='*.go'` returns zero results (excluding justfile-canonical-e2e which is Task 2)
- Duplicate TestMain in `test-knowledge-convention-driven/` removed; that package uses the parent's `forgeBinary`
- Walk-up path resolution (`forge-cli/bin/forge`) eliminated from all test files

## Hard Rules
- Do not change test assertions or test logic
- The `tests/e2e/features/test-knowledge-convention-driven/` sub-package must import and use the parent package's `forgeBinary` via an exported var or test helper

## Implementation Notes
- Some feature tests may be in sub-packages (different Go packages but same module) — they need to import the parent package's exported variable
- Use `grep -rn 'exec.Command("forge"' tests/e2e/ --include='*.go'` to find all instances before fixing
