---
name: fix-bug
description: Systematically fix a bug using TDD workflow — reproduce, write failing tests, fix, verify. Ensures the bug is captured by tests before any code changes.
allowed-tools: Bash Read Write Edit Grep Glob Agent LSP
argument-hint: "[error-msg] [scope]"
---

# /fix-bug

Systematic bug fix workflow: **Reproduce → Test → Fix → Verify**.

Core principle: never touch production code until a failing test proves the bug exists.

<CODING_PRINCIPLES>
- Think Before Coding: Before writing any fix, restate the bug and its suspected root cause in your own words. Verify your diagnosis against the evidence — do not jump to the first plausible explanation. If the root cause is ambiguous, investigate further before proceeding.
- Simplicity First: Fix only what is broken. No speculative changes, no "while I'm here" improvements, no refactoring. Trivial fixes (typos, config) use judgment — full analysis is not needed.
- Surgical Changes: Modify only the code directly on the failing path. Do not touch neighboring code, reformat unrelated lines, or refactor tangential logic. Scope boundary = failing code path only.
- Goal-Driven Execution: Define a clear, verifiable success condition: the specific test(s) that fail before the fix and pass after. After implementation, confirm the condition is met — if not, iterate.
</CODING_PRINCIPLES>

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
- Tests and fix must be committed together in a single atomic commit
- If the bug cannot be reproduced (Step 2), STOP and report — do not write tests or fix code for an unconfirmed bug
- All quality gate checks must pass before committing (Step 5)
</EXTREMELY-IMPORTANT>

---

## Step 1: Understand the Bug

Collect all available context before touching any code.

**Project Knowledge**: Infer relevant domains from bug description, affected files, and scope.
Check `docs/conventions/` and `docs/business-rules/` for project-specific knowledge relevant to this task.
Read each file's YAML frontmatter `domains` field to determine relevance.
Load files whose domains overlap with the task context.
If no files match, skip — no matching convention files for this task.

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

**Guidelines:**
- Follow <CODING_PRINCIPLES> — Simplicity First and Surgical Changes define the scope boundary
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

Strict sequential order. Stop at first failure. See Forge Guide Quality Gate Protocol for failure actions.

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
Skill(skill="forge:git-commit")
```

The commit must include both the fix and the tests in a single atomic commit.

Commit message format:
```
fix(<scope>): <what was wrong and is now correct>

Root cause: <one sentence>
Fixes: #<issue-number> (if applicable)
```

---

## Knowledge Review

After Step 6 (Commit) completes, run knowledge auto-extraction from the bug fix session.

### Parameters

| Parameter | Value |
|-----------|-------|
| `trigger` | `fix-bug` |
| `artifacts` | root cause analysis (the "Root cause: \<why\>" note from Step 4), fix approach (files changed and the nature of the fix) |

### Artifact Scanning Scope

Focus on non-obvious root causes and debugging patterns.

### Knowledge Types

The extraction routine identifies four knowledge types:

| Type | Target | Format reference |
|------|--------|-----------------|
| Decision | `docs/decisions/<type>.md` | `decision-logging.md` Section 6 (row format), Section 7 (manifest update) |
| Lesson | `docs/lessons/<slug>.md` | `learn/templates/lesson-entry.md` |
| Convention | `docs/conventions/<topic>.md` | `/consolidate-specs` tech-specs entry format, with project-global ID |
| Business Rule | `docs/business-rules/<domain>.md` | `/consolidate-specs` biz-specs entry format, with project-global ID |

### Extraction Flow

#### Step 1: Scan artifacts

Read all artifacts specified above.

#### Step 2: Identify notable knowledge

Apply the "notable knowledge" heuristics below to determine if any notable knowledge exists in the scanned artifacts. Classify each candidate by knowledge type (Decision, Lesson, Convention, Business Rule). Filter out trivial fixes (typos, simple config changes, obvious mistakes).

#### Step 3: Vocabulary-assisted classification

If `/consolidate-specs` has previously generated vocabulary (from drift-detection runs), use the domain keywords from existing `docs/conventions/` and `docs/business-rules/` files to suggest which target file each extracted item belongs to. This is a suggestion — the agent makes the final classification decision based on content.

If no vocabulary exists (no prior `/consolidate-specs` run), classify unassisted using the domain-to-file mapping from `/consolidate-specs` skill Step 5.

#### Step 4: Silent exit if no notable knowledge

If no candidates pass the "notable" heuristics (below), **produce no output**. Do not ask the user anything. Return silently.

#### Step 5: Present for user confirmation

Use AskUserQuestion to present extracted candidates:

```
Knowledge extracted from fix-bug:

  [1] <Decision> → docs/decisions/<type>.md
  [2] <Lesson> → docs/lessons/<slug>.md
  [3] <Convention> → docs/conventions/<topic>.md
  [4] <Business Rule> → docs/business-rules/<domain>.md

Enter numbers to save (comma-separated), or all / none:
```

User input handling:
- `none` → discard all candidates, no output
- `all` → save all candidates
- comma-separated numbers → save only selected candidates

#### Step 6: Write confirmed knowledge

For each confirmed candidate, write to the target file using the format defined by the knowledge type. Create target files if they do not exist. When creating new convention/business-rule files, include YAML frontmatter with `title` and `domains` per `/consolidate-specs` Domain Derivation Rules.

Do NOT write to knowledge directories without explicit user confirmation from Step 5.

### Notable Knowledge Heuristics

The heuristics determine whether a piece of knowledge is "notable" (worth extracting) vs "routine" (skip silently). The goal is a false-positive rate below 30%.

**Decisions — NOT notable when:**

- The choice is the standard/default option in the ecosystem (e.g., "used standard library", "used ORM for database access")
- No meaningful alternatives existed (e.g., "used the only available API")
- The decision is purely cosmetic or stylistic with no architectural impact
- The decision replicates an existing entry in `docs/decisions/`

**Decisions — NOTABLE when:**

- Multiple viable alternatives existed and the choice has lasting impact (e.g., "chose event-driven over polling for state sync")
- The decision involves a non-obvious tradeoff (e.g., "sacrificed consistency for availability in the cache layer")
- A constraint forced an unconventional approach (e.g., "used file-based locking because the Redis dependency was disallowed")

**Lessons — NOT notable when:**

- The root cause is a trivial mistake (e.g., typo, missing import, wrong variable name)
- The issue is standard to the framework/language (e.g., "null pointer from uninitialized field")
- The fix was obvious from the error message
- The lesson replicates an existing entry in `docs/lessons/`

**Lessons — NOTABLE when:**

- The root cause was non-obvious (e.g., race condition from hidden shared state, ordering dependency across services)
- The debugging path was indirect (e.g., "symptom appeared in module A but root cause was in module B")
- The issue would recur in similar contexts and the pattern is worth documenting (e.g., "non-thread-safe map in concurrent handler")

**Conventions — NOT notable when:**

- The pattern is already documented in `docs/conventions/`
- The pattern is a one-off choice specific to this feature
- The pattern is standard practice in the ecosystem (e.g., "used REST for HTTP API")

**Conventions — NOTABLE when:**

- The pattern should be repeated across the project (e.g., "all CLI commands use cobra with this flag structure")
- A project-specific standard was established (e.g., "config files use YAML with this schema structure")
- The pattern emerged from implementation and was not pre-designed

**Business Rules — NOT notable when:**

- The rule is feature-specific logic (e.g., "this feature's form validates email format")
- The rule is a standard CRUD constraint (e.g., "required fields must be non-empty")
- The rule replicates an existing entry in `docs/business-rules/`

**Business Rules — NOTABLE when:**

- The rule applies across features (e.g., "all monetary values use integer cents, never float")
- The rule expresses a domain invariant (e.g., "order status can only advance, never regress")
- The rule constrains user-facing behavior across the system (e.g., "all user actions require authentication except health-check endpoints")

### Deduplication

Before presenting candidates in Step 5, check for duplicates:

1. **Decisions**: grep `docs/decisions/<type>.md` for similar Decision text
2. **Lessons**: grep `docs/lessons/*.md` for similar Root Cause content
3. **Conventions**: grep `docs/conventions/*.md` for similar rule descriptions
4. **Business Rules**: grep `docs/business-rules/*.md` for similar rule statements

If a duplicate is found, exclude the candidate and do not present it. The heuristic goal is: if it is already documented, do not re-extract it.

### Rules

- Extraction logic must be **conservative**: only extract genuinely non-obvious knowledge
- Must not write to knowledge directories without explicit user confirmation
- Silent when no notable knowledge is detected — no output, no prompts
- All output formats must be compatible with `/learn` skill and `/consolidate-specs` overlap detection
- Deduplication runs before presentation — never present a duplicate of existing knowledge

---

## Output Summary

After completion (and optional knowledge review), report:

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
