---
created: "2026-05-16"
tags: [testing, architecture]
---

# E2E Test Generation Produces Low-Quality Tests Without Quality Gates

## Problem

67 e2e test cases were audited; 35 (52%) had quality issues:

| Antipattern | Count | Impact |
|---|---|---|
| Recursive `go test ./...` inside test | 2 | Process explosion (126+ orphaned processes) |
| Duplicate tests in root + features/ | 12 | Wasted CI time, maintenance burden |
| Unconditional `t.Skip` (never implemented) | 2 | Dead code, false coverage signal |
| Conditional skip without self-contained fixture | 19 | Tests silently pass or skip depending on environment state |

## Root Cause

Causal chain (4 levels):

1. **Symptom**: 52% of e2e tests are ineffective or actively harmful.
2. **Direct cause**: Each test category has a specific generation flaw:
   - **Recursive**: Meta-tests verify "the test suite compiles/passes" by invoking `go test ./...` from within the suite itself.
   - **Duplicate**: The graduation workflow (`/graduate-tests`) copies tests from `tests/e2e/features/<slug>/` to `tests/e2e/` root without removing the source. Both are separate Go packages under the same module, so `go test ./...` runs them twice.
   - **Skip**: Agent-generated placeholders for scenarios requiring complex environment setup (git worktree, feature branch) that were never implemented.
   - **Conditional skip**: Tests rely on live project state (e.g., pending tasks exist) instead of creating isolated fixtures.
3. **Deeper cause**: The forge test generation pipeline (`/gen-test-cases` → `/gen-test-scripts`) optimizes for coverage quantity, not test quality. There is no post-generation quality gate that checks:
   - Can this test actually run in isolation?
   - Does this test assert anything meaningful, or does it have vacuous truth paths?
   - Does this test duplicate an existing test?
   - Does this test recursively invoke itself?
4. **Fundamental cause**: Agent-generated code lacks the "smell test" that human reviewers apply. A human would immediately flag `t.Skip` without a plan to implement, or `go test ./...` inside a test function, or tests without self-contained setup. The generation pipeline has no equivalent validation step.

## Solution

### Immediate fixes per antipattern

**Recursive tests**: Add environment variable recursion guard (see `gotcha-recursive-go-test-process-explosion.md`).

**Duplicate tests**: After graduation, delete the source file from `tests/e2e/features/<slug>/`. Or: don't graduate — keep tests only in features/ subdirectories.

**Skip placeholders**: Either implement the test (with proper fixture setup) or delete it. A `t.Skip` that's never resolved is a false coverage signal.

**Conditional skip without fixture**: Refactor to create self-contained fixtures. Every test must set up its own world:

```go
// BAD: depends on environment having pending tasks
func TestTC_001(t *testing.T) {
    out, code := runCLI("forge", "task", "claim")
    if code != 0 {
        t.Skip("no pending tasks")
    }
    // assertions...
}

// GOOD: self-contained fixture
func TestTC_001(t *testing.T) {
    dir := setupProjectWithPendingTasks(t)  // creates temp dir + task files
    out, code := runCLIInDir(dir, "forge", "task", "claim")
    require.Equal(t, 0, code)
    // assertions...
}
```

### Structural fix: test quality gate

Add a validation step to `/gen-test-scripts` (or `/eval-test-cases`) that rejects generated tests with:

1. **Recursion check**: Grep for `exec.Command("go", "test"` or equivalent in the generated test. If found, require a recursion guard.
2. **Skip check**: Reject `t.Skip` without a linked TODO issue or environment-detection rationale.
3. **Fixture check**: Every test that calls an external CLI must use `t.TempDir()` and set up its own project structure.
4. **Duplicate check**: Before generating, scan existing test files for matching `func TestTC_*` names.
5. **Vacuous assertion check**: Reject `if condition { assert.X(...) }` patterns — the assertion must always execute, or the test must skip explicitly when the condition is unmet.

## Reusable Pattern

**Test quality is not test quantity.** When generating tests automatically (via agent or script), add a post-generation validation gate that checks for these antipatterns. The generation step and the validation step should be separate — generating and validating in the same agent pass means the agent won't question its own output.

Specifically for forge's `/gen-test-scripts` skill: the generated scripts should be validated by `/eval-test-cases` against a rubric that includes these antipattern checks. Currently the rubric evaluates test case completeness but not test code quality.

## Related Files

- `tests/e2e/simplify_e2e_tests_cli_test.go` (recursive tests)
- `tests/e2e/cli_lean_output_cli_test.go` (conditional skip, no fixtures)
- `tests/e2e/cli_list_reverse_chronological_cli_test.go` (duplicate of features/ version)
- `tests/e2e/fix_task_claim_priority_cli_test.go` (duplicate of features/ version)
- `tests/e2e/feature_set_command_cli_test.go` (skip placeholders at TC-016, TC-017)
- `docs/lessons/gotcha-recursive-go-test-process-explosion.md` (related lesson on recursion)
