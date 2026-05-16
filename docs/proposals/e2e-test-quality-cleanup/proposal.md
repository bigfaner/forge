---
created: "2026-05-16"
author: "faner"
status: Draft
---

# Proposal: E2E Test Quality Cleanup & Forge Pipeline Hardening

## Problem

52% of e2e tests (35/67) have quality antipatterns — some actively dangerous (recursive process explosion), others silently harmful (dead tests giving false coverage signals). The forge test generation pipeline has no quality gate to prevent regenerating these antipatterns.

### Evidence

Audit of `tests/e2e/` revealed:

| Antipattern | Count | Impact |
|---|---|---|
| Recursive `go test ./...` inside test | 2 | Process explosion (126+ orphaned processes on Windows) |
| Duplicate tests in root + features/ | 12 | Wasted CI time, maintenance burden |
| Unconditional `t.Skip` (never implemented) | 2 | Dead code, false coverage signal |
| Conditional skip without self-contained fixture | 19 | Tests silently pass or skip depending on environment |

Source: `docs/lessons/gotcha-recursive-go-test-process-explosion.md` and `docs/lessons/gotcha-e2e-test-quality-antipatterns.md`

### Urgency

The recursive tests can consume 6GB+ RAM and make the machine unresponsive. Dead tests mask real coverage gaps. Without a pipeline fix, `/gen-test-scripts` will keep producing the same antipatterns.

## Proposed Solution

Two-phase approach:

1. **Phase 1 — Clean up existing tests**: Fix recursion guards, handle duplicate tests, delete dead placeholders, refactor conditional skips to self-contained fixtures.
2. **Phase 2 — Harden forge pipeline**: Extend `test-cases` rubric with antipattern detection, fix `graduate-tests` to delete source files after copy.

## Requirements Analysis

### Key Scenarios

- Developer runs `just test-e2e` — no process explosion, no duplicate test execution
- Agent generates test scripts via `/gen-test-scripts` — rubric catches antipatterns before merge
- `/graduate-tests` runs — source files cleaned up, no duplicates remain
- CI runs e2e suite — all tests either pass or fail with meaningful assertions (no silent skips)

### Non-Functional Requirements

- All e2e tests must run deterministically regardless of environment state
- Test execution time should not increase (removing duplicates will decrease it)

### Constraints & Dependencies

- Go test framework on Windows (no `PR_SET_PDEATHSIG`, child processes survive parent death)
- Existing `tests/e2e/features/<slug>/` structure must be preserved

## Alternatives & Industry Benchmarking

### Comparison Table

| Approach | Pros | Cons | Verdict |
|----------|------|------|---------|
| Do nothing | Zero effort | Recursion risk persists, coverage signal unreliable | Rejected: active danger |
| Delete all bad tests, no pipeline fix | Fast cleanup | Antipatterns will reappear on next generation | Rejected: treats symptom |
| Separate lint script for antipatterns | Strong enforcement | Adds a new tool to maintain; separate from existing eval flow | Rejected: overengineered for this scope |
| **Fix tests + extend rubric** | Fixes root cause, unified in existing eval flow | Slightly larger scope | **Selected: comprehensive but proportional** |

## Scope

### In Scope

- Add recursion guard to `simplify_e2e_tests_cli_test.go` TC-003/TC-004
- Delete dead `t.Skip` placeholders TC-016/TC-017 from `feature_set_command_cli_test.go`
- Refactor `cli_lean_output_cli_test.go` to use self-contained fixtures
- Fix `graduate-tests` skill to delete source files after copying
- Extend `rubrics/test-cases.md` with antipattern detection dimension
- Update `gen-test-scripts` SKILL.md to reference the new rubric dimension

### Out of Scope

- Rewriting the entire e2e test framework
- Adding a new standalone lint tool
- Changing the test profile system
- Fixing `forge-info-commands` feature (BLOCKED status — separate concern)

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Self-contained fixtures for cli_lean_output are complex (need forge project setup) | M | M | Create a shared helper that sets up temp project with tasks/features via `t.TempDir()` |
| Rubric extension changes eval scores for existing test cases | L | L | Add as new dimension, don't redistribute points from existing ones |
| graduate-tests source deletion breaks features with cross-feature dependencies | L | M | Only delete the specific feature's source files, not the entire features/ directory |

## Success Criteria

- [ ] `simplify_e2e_tests` TC-003/TC-004 have `FORGE_E2E_RECURSION_GUARD` check; no process explosion on Windows
- [ ] Zero unconditional `t.Skip` in `tests/e2e/` (TC-016/TC-017 deleted)
- [ ] `cli_lean_output` tests run deterministically without depending on live forge state
- [ ] `graduate-tests` deletes source files from `features/<slug>/` after successful copy
- [ ] `rubrics/test-cases.md` includes antipattern checks (recursion, skip, fixture, duplicate, vacuous assertion)
- [ ] `just test-e2e` completes without orphans, duplicates, or unexpected skips

## Next Steps

- Proceed to `/quick-tasks` to generate task list
