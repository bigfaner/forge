---
name: fix-bug
description: Systematically fix a bug using TDD workflow — reproduce, write failing tests, fix, verify. Ensures the bug is captured by tests before any code changes.
allowed-tools: Bash Read Write Edit Grep Glob Agent LSP AskUserQuestion
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

**Argument parsing**: The `argument-hint` is `[error-msg] [scope]`. When the user invokes `/fix-bug <args>`, parse `<args>` as follows:
- If `<args>` is empty: `--issue` is derived from Step 1 investigation (user provides details interactively)
- If `<args>` contains a quoted string or unquoted text: that text becomes `--issue` (bug description)
- If `<args>` contains a path-like segment (contains `/` or `.\`): that segment becomes `--scope`
- If a single path-like segment is provided with no issue description: `--scope` is set, `--issue` is derived interactively
- Examples: `/fix-bug "login throws 500"` → `--issue "login throws 500"`; `/fix-bug src/auth/` → `--scope src/auth/`; `/fix-bug "null pointer" pkg/handler/` → both set

## Workflow

```
1. Understand → 2. Reproduce → 3. Write failing tests → 4. Fix → 5. Verify → 6. Commit
```

<EXTREMELY-IMPORTANT>
- Never touch production code until a failing test proves the bug exists (Step 3 before Step 4)
- Tests and fix must be committed together in a single atomic commit
- All quality gate checks must pass before committing (Step 5)
</EXTREMELY-IMPORTANT>

---

## Step 1: Understand the Bug

Collect all available context before touching any code.

**Project Knowledge**: Infer relevant domains from bug description, affected files, and scope.
Check `docs/conventions/` and `docs/business-rules/` for project-specific knowledge relevant to this task.
Read each file's YAML frontmatter `domains` field to determine relevance.
Load files whose domains match the affected scope or keywords from the bug description.
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
# Run just unit-test to establish baseline
just unit-test [scope]
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

Run `just unit-test` (pass `[scope]` from `--scope` parameter if specified) — it **must fail** before the fix:

```bash
just unit-test [scope]
```

<HARD-RULE>
If the new unit test passes before any fix, the test does not capture the bug. Revise the test until it fails for the right reason.
</HARD-RULE>

### 3b. E2E / Integration Test (skip if `--skip-e2e`)

Add an e2e test only when the bug is observable at the API, CLI, or UI surface.

| Bug surface | Test location | Runner |
|-------------|--------------|--------|
| UI behavior | `tests/<journey>/ui.spec.ts` | 浏览器自动化 (test profile) |
| API endpoint | `tests/<journey>/api.spec.ts` | HTTP 客户端 |
| CLI command | `tests/<journey>/cli.spec.ts` | 子进程执行 |
| Mobile | `tests/<journey>/mobile.yaml` | Maestro YAML |
| TUI | `tests/<journey>/tui.spec.ts` | 子进程 + stdin pipe |

Bug fix tests go to the journey directory corresponding to the affected surface.

Run `just test` — it **must fail** before the fix:

```bash
just test
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

Execute the quality gate sequence (pass `[scope]` from `--scope` parameter if specified):

```
just compile [scope] → just fmt [scope] → just lint [scope] → just unit-test [scope]
```

Strict sequential order. Stop at first failure. On failure: compile → fix & retry; fmt → non-blocking warning; lint → self-fix (1 retry) then blocked; unit-test → fix & retry.

Surface-level tests (if written in Step 3b):

```bash
just test
```

**Verification checklist:**
- [ ] New unit test(s): PASS
- [ ] New surface-level test(s): PASS (if written)
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

| Type | Target |
|------|--------|
| Decision | `docs/decisions/<type>.md` |
| Lesson | `docs/lessons/<slug>.md` |
| Convention | `docs/conventions/<topic>.md` |
| Business Rule | `docs/business-rules/<domain>.md` |

<!-- INLINE:origin=learn/templates/ -->
### Write Formats (per knowledge type)

**Decision** — append one table row to `docs/decisions/<type>.md`:
```
| {{DATE}} | {{FEATURE_SLUG}} | {{DECISION}} | {{RATIONALE}} | {{SOURCE}} |
```
Field constraints: Date=ISO 8601, Feature=slug or `-`, Decision/Rationale=max 80 chars each, Source=file path or `manual`. If type file missing, create with standard header (`# {{TYPE_NAME}} Decisions` + column headers). After append, update `docs/decisions/manifest.md` categories count and insert into Recent Decisions (newest first, max 10 rows).

**Lesson** — create `docs/lessons/<prefix>-<slug>.md` with sections: Problem, Root Cause, Solution, Reusable Pattern, Example (optional), Related Files (optional), References (optional). Tags: select 1-4 from fixed vocabulary (architecture, interface, data-model, dependencies, error-handling, testing, security, local-dev-deployment). File prefix by category: `debug-`, `arch-`, `tool-`, `pattern-`, `gotcha-`.

**Convention** — append to `docs/conventions/<topic>.md` using project-global ID `TECH-{{TOPIC}}-{{NNN}}`:
```markdown
### TECH-{{TOPIC}}-{{NNN}}: {{SPEC_TITLE}}
**Requirement**: {{CONCISE_REQUIREMENT}}
**Scope**: [CROSS]
**Source**: /learn entry {{DATE}}
```
Sequence: max existing NNN in target file + 1.

**Business Rule** — append to `docs/business-rules/<domain>.md` using project-global ID `BIZ-{{DOMAIN}}-{{NNN}}`:
```markdown
### BIZ-{{DOMAIN}}-{{NNN}}: {{RULE_TITLE}}
**Rule**: {{CONCISE_RULE_STATEMENT}}
**Context**: {{WHY_THIS_RULE_EXISTS}}
**Scope**: [CROSS]
**Source**: /learn entry {{DATE}}
```
For new convention/business-rule files, include YAML frontmatter with `title` and `domains` (3-7 specific keywords derived from entry content).

<!-- INLINE:origin=consolidate-specs/rules/ -->
### Classification & ID Rules

- **CROSS**: referenced by 2+ features, expresses domain invariant, or establishes naming/error-handling convention
- **LOCAL**: only meaningful within this feature's scope
- **Project-global ID**: `BIZ-<domain>-<NNN>` or `TECH-<topic>-<NNN>` — prefix from target filename, sequence = max existing NNN + 1

### Domain-to-Decision Mapping (for vocabulary-assisted classification)

| Spec domain keywords | Decision file |
|---------------------|---------------|
| system structure, layering, modules, architecture | `architecture.md` |
| API contracts, data shapes, serialization, interface | `interface.md` |
| schema, indexing, soft-delete, data model | `data-model.md` |
| libraries, versions, packages, dependencies | `dependencies.md` |
| error types, status codes, error propagation | `error-handling.md` |
| test patterns, coverage, mocking | `testing.md` |
| auth, permissions, data protection, security | `security.md` |
| dev environment, tooling, deployment | `local-dev-deployment.md` |
| naming, conventions, coding standards | `architecture.md` |
| validation, state transitions, calculation | closest match or `architecture.md` |

### Extraction Flow

1. **Scan artifacts** — read all artifacts specified in Parameters above.
2. **Identify notable knowledge** — apply heuristics below. Classify by type. Filter out trivial fixes.
3. **Vocabulary-assisted classification** — use existing `docs/conventions/` and `docs/business-rules/` domains frontmatter to suggest target files; if none exist, use the Domain-to-Decision Mapping table above.
4. **Silent exit if no notable knowledge** — produce no output, return silently.

5. **Auto-save configuration check** — run `forge config get auto.knowledgeSave`:

   | Exit Code | Mode value | Action |
   |-----------|-----------|--------|
   | 0 | `true` | Skip step 6. Treat all candidates as confirmed, proceed directly to step 7. |
   | 0 | `false` | Present step 6 confirmation. |
   | Non-zero | — | Fallback: present step 6 confirmation. |

   Mode context: `quick` via `/quick` pipeline, `full` via full pipeline. Parse the config output format `quick:<val> full:<val>` (e.g., `quick:true full:false`) and select the value matching the current mode.

6. **Present for user confirmation** (skipped when auto-save enabled) — use AskUserQuestion:

   ```
   Knowledge extracted from fix-bug:

     [1] <Decision> → docs/decisions/<type>.md
     [2] <Lesson> → docs/lessons/<slug>.md
     [3] <Convention> → docs/conventions/<topic>.md
     [4] <Business Rule> → docs/business-rules/<domain>.md

   Enter numbers to save (comma-separated), or all / none:
   ```

   User input: `none` → discard; `all` → save all; comma-separated numbers → save selected only.

7. **Write confirmed knowledge** — write each candidate to its target file using the Write Formats above. Create files if needed. For new convention/business-rule files, include YAML frontmatter with `title` and `domains`. Do NOT write without explicit user confirmation from step 6.

### Notable Knowledge Heuristics

Goal: false-positive rate below 30%. Only extract genuinely non-obvious knowledge.

| Type | Skip (routine) | Extract (notable) |
|------|---------------|-------------------|
| **Decision** | Standard/default option, no alternatives, cosmetic, or duplicate | Multiple viable alternatives with lasting impact, non-obvious tradeoff, constraint-forced unconventional approach |
| **Lesson** | Trivial mistake, framework-standard issue, obvious from error message, or duplicate | Non-obvious root cause (race condition, cross-module), indirect debugging path, recurring pattern worth documenting |
| **Convention** | Already documented, one-off choice, ecosystem-standard, or feature-specific | Should be repeated across project, project-specific standard established, emerged from implementation |
| **Business Rule** | Feature-specific logic, standard CRUD constraint, or duplicate | Applies across features, expresses domain invariant, constrains user-facing behavior system-wide |

### Deduplication

Before presenting candidates in Step 5, grep target directories for similar existing entries. Exclude any candidate that duplicates already-documented knowledge:
- **Decisions**: grep `docs/decisions/<type>.md` for similar Decision text
- **Lessons**: grep `docs/lessons/*.md` for similar Root Cause content
- **Conventions**: grep `docs/conventions/*.md` for similar rule descriptions
- **Business Rules**: grep `docs/business-rules/*.md` for similar rule statements

### Rules

- Extraction logic must be **conservative**: only extract genuinely non-obvious knowledge
- Must not write to knowledge directories without explicit user confirmation
- Silent when no notable knowledge is detected — no output, no prompts
- All output formats must follow the Write Formats and Classification & ID Rules sections above
- Deduplication runs before presentation — never present a duplicate of existing knowledge

<!-- END INLINE:origin=consolidate-specs/rules/ -->

---

## Output Summary

After completion (and optional knowledge review), report:

```
Bug Fix Summary
───────────────
Bug:     <description>
Fix:     <file(s) changed>
Tests:   <N unit tests added> + <M surface-level tests added>
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
| Skipping e2e when the bug is user-facing | Add at least one surface-level smoke test |
| Committing fix and tests separately | One atomic commit: fix + tests together |
