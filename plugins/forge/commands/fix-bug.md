---
name: fix-bug
description: Systematically fix a bug using TDD workflow — reproduce, write failing tests, fix, verify. Ensures the bug is captured by tests before any code changes.
allowed_tools: ["Bash", "Read", "Write", "Edit", "Grep", "Glob", "Agent", "LSP"]
argument-hints:
  - name: error-msg
    description: Error message, stack trace, or symptom description to locate the bug (e.g. "TypeError: Cannot read property 'id' of undefined")
    required: false
  - name: scope
    description: Affected module or package path to narrow the search (e.g. src/parser, pkg/auth). Auto-detected if omitted.
    required: false
---

# /fix-bug

Systematic bug fix workflow: **Reproduce → Test → Fix → Verify**.

Core principle: never touch production code until a failing test proves the bug exists.

## Parameters

| Parameter | Default | Description |
|-----------|---------|-------------|
| `--issue` | — | GitHub issue number or bug description |
| `--scope` | auto | Affected module/package path (auto-detected if omitted) |
| `--skip-e2e` | false | Skip e2e tests (use when no UI/API surface is affected) |

## Workflow

```
1. Understand → 2. Reproduce → 3. Write failing tests → 4. Fix → 5. Verify → 6. Commit
```

<EXTREMELY-IMPORTANT>
- Never touch production code until a failing test proves the bug exists (Step 3 before Step 4)
- Fix only what the failing tests require — no scope creep, no refactoring, no "improvements"
- Tests and fix must be committed together in a single atomic commit
- If the bug cannot be reproduced (Step 2), STOP and report — do not write tests or fix code for an unconfirmed bug
- All quality gate checks must pass before committing (Step 5)
</EXTREMELY-IMPORTANT>

---

## Step 1: Understand the Bug

Collect all available context before touching any code.

**Project Knowledge**: Read relevant project knowledge files:
- Infer relevant domains from bug description, affected files, and scope
- Read matching files from `docs/business-rules/` and `docs/conventions/`
- Example mappings: "auth"/"login"/"permission" → `business-rules/auth.md`; "state"/"validation"/"lifecycle" → `business-rules/<domain>.md`; "API"/"endpoint"/"route" → `conventions/api.md`; "error"/"status code" → `conventions/error-handling.md`; "database"/"schema"/"migration" → `conventions/data-model.md`; "test"/"mock"/"coverage" → `conventions/testing.md`
- If no matching file exists, skip this step

**Gather:**
- Bug description / error message / stack trace
- Steps to reproduce (from issue or user)
- Expected vs. actual behavior
- Affected version / commit (use `git log --oneline -20` if unknown)

**Locate the blast radius:**

```bash
# Find files related to the symptom
grep -r "<error-keyword>" src/ --include="*.ts" -l
git log --oneline --all -- <suspected-file>
```

Write a one-paragraph **Bug Summary** before proceeding:

```
Bug: <what goes wrong>
Trigger: <exact steps or input that causes it>
Expected: <correct behavior>
Actual: <observed behavior>
Suspected location: <file:line or module>
```

---

## Step 2: Reproduce

Confirm the bug is reproducible in the current codebase before writing any tests.

```bash
# Run just test to establish baseline
just test [scope]
```

**Reproduction checklist:**
- [ ] Bug is reproducible on current branch
- [ ] Baseline test suite passes (no pre-existing failures masking the bug)
- [ ] Exact reproduction steps documented

<HARD-GATE>
If the bug cannot be reproduced, STOP. Report to the user with findings. Do not proceed to write tests or fix code for an unconfirmed bug.
</HARD-GATE>

---

## Step 3: Write Failing Tests

Write tests that **fail because of the bug** and will **pass after the fix**. This is the proof that the fix works.

### 3a. Unit Test

Locate the test file closest to the buggy code. Add a focused test case:

```
Convention: test file lives next to source file
  src/foo/bar.ts  →  src/foo/bar.test.ts  (or bar.spec.ts)
  pkg/foo/bar.go  →  pkg/foo/bar_test.go
  foo/bar.py      →  tests/test_bar.py
```

Test naming convention:
```
"bug: <short description of the incorrect behavior>"
// e.g. "bug: returns null when input is empty string"
```

Run `just test [scope]` — it **must fail** before the fix (apply **Scope Resolution** from Forge Guide before invoking):

```bash
just test [scope]
```

<HARD-RULE>
If the new unit test passes before any fix, the test does not capture the bug. Revise the test until it fails for the right reason.
</HARD-RULE>

### 3b. E2E / Integration Test (skip if `--skip-e2e`)

Add an e2e test only when the bug is observable at the API, CLI, or UI surface.

| Bug surface | Test location | Runner |
|-------------|--------------|--------|
| UI behavior | `tests/e2e/features/<slug>/ui.spec.ts` | Playwright |
| API endpoint | `tests/e2e/features/<slug>/api.spec.ts` | fetch |
| CLI command | `tests/e2e/features/<slug>/cli.spec.ts` | child_process |

Bug fix tests go to the `features/` staging area, same as feature tests. This ensures `just test-e2e --feature <slug>` can discover and run them.

Run `just test-e2e --feature <slug>` — it **must fail** before the fix:

```bash
just test-e2e --feature <slug>
```

---

## Step 4: Fix

With failing tests in place, implement the minimal fix.

**Principles:**
- Fix only what the failing tests require — no scope creep
- Do not refactor surrounding code
- Do not add features or "improvements"
- If the root cause is in a dependency, document it and apply the minimal workaround

**Root cause note** — before moving on, write one sentence in a code comment or commit body:
```
// Root cause: <why this happened, not just what changed>
```

---

## Step 5: Verify (Quality Gate)

Execute the quality gate sequence. Apply **Scope Resolution** from the Forge Guide for each command:

```
just compile [scope] → just fmt [scope] → just lint [scope] → just test [scope]
```

Strict sequential order. Stop at first failure:

| Failed step | Action |
|---|---|
| `compile` | Fix compilation errors, then retry from compile |
| `fmt` | Mark task as `blocked` (auto-fix failed = toolchain issue) |
| `lint` | Self-fix (max 1 retry), then mark `blocked` if still failing |
| `test` | Fix failing tests, then retry from compile |

E2E (if written in Step 3b):

```bash
just test-e2e --feature <slug>
```

**Verification checklist:**
- [ ] New unit test(s): PASS
- [ ] New e2e test(s): PASS (if written)
- [ ] Pre-existing tests: no regressions
- [ ] Build succeeds (if applicable)

<HARD-GATE>
Do not proceed to commit if any pre-existing test is newly failing. Investigate whether the fix introduced a regression.
</HARD-GATE>

---

## Step 6: Commit

```
Skill(skill="git-commit")
```

The commit must include both the fix and the tests in a single atomic commit.

Commit message format:
```
fix(<scope>): <what was wrong and is now correct>

Root cause: <one sentence>
Fixes: #<issue-number> (if applicable)
```

---

## Output Summary

After completion, report:

```
Bug Fix Summary
───────────────
Bug:     <description>
Fix:     <file(s) changed>
Tests:   <N unit tests added> + <M e2e tests added>
Result:  All tests pass ✓
Commit:  <commit hash>
```

---

## Common Pitfalls

| Pitfall | Correct approach |
|---------|-----------------|
| Fixing before writing a failing test | Always write the test first |
| Test passes before fix | Test doesn't capture the bug — revise it |
| Fixing more than the bug | Minimal fix only; open a separate task for cleanup |
| Skipping e2e when the bug is user-facing | Add at least one e2e smoke test |
| Committing fix and tests separately | One atomic commit: fix + tests together |
